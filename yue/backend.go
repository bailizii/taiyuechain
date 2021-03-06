// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package yue implements the Taiyuechain protocol.
package yue

import (
	"errors"
	"fmt"
	"math/big"
	"runtime"
	"strconv"
	"sync"

	"github.com/taiyuechain/taiyuechain/consensus/tbft"
	"github.com/taiyuechain/taiyuechain/crypto"
	config "github.com/taiyuechain/taiyuechain/params"

	//"sync/atomic"

	"github.com/taiyuechain/taiyuechain/common"
	"github.com/taiyuechain/taiyuechain/common/hexutil"

	"github.com/taiyuechain/taiyuechain/accounts"
	"github.com/taiyuechain/taiyuechain/consensus"
	elect "github.com/taiyuechain/taiyuechain/consensus/election"
	ethash "github.com/taiyuechain/taiyuechain/consensus/minerva"
	"github.com/taiyuechain/taiyuechain/core"
	"github.com/taiyuechain/taiyuechain/core/bloombits"
	"github.com/taiyuechain/taiyuechain/core/types"
	"github.com/taiyuechain/taiyuechain/core/vm"
	"github.com/taiyuechain/taiyuechain/log"
	"github.com/taiyuechain/taiyuechain/rlp"

	"github.com/taiyuechain/taiyuechain/event"
	"github.com/taiyuechain/taiyuechain/internal/taiapi"
	"github.com/taiyuechain/taiyuechain/yue/downloader"
	"github.com/taiyuechain/taiyuechain/yue/filters"
	"github.com/taiyuechain/taiyuechain/yue/gasprice"
	"github.com/taiyuechain/taiyuechain/yuedb"

	//"github.com/taiyuechain/taiyuechain/miner"
	"github.com/taiyuechain/taiyuechain/cim"
	"github.com/taiyuechain/taiyuechain/node"
	"github.com/taiyuechain/taiyuechain/p2p"
	"github.com/taiyuechain/taiyuechain/params"
	"github.com/taiyuechain/taiyuechain/rpc"
)

type LesServer interface {
	Start(srvr *p2p.Server)
	Stop()
	Protocols() []p2p.Protocol
	SetBloomBitsIndexer(bbIndexer *core.ChainIndexer)
}

// Taiyuechain implements the Taiyuechain full node service.
type Taiyuechain struct {
	config      *Config
	chainConfig *params.ChainConfig

	// Channel for shutting down the service
	shutdownChan chan bool // Channel for shutting down the Taiyuechain

	// Handlers
	txPool *core.TxPool

	//snailPool *chain.SnailPool

	agent    *PbftAgent
	election *elect.Election

	blockchain *core.BlockChain
	//snailblockchain *chain.SnailBlockChain
	protocolManager *ProtocolManager
	lesServer       LesServer

	// DB interfaces
	chainDb yuedb.Database // Block chain database

	eventMux       *event.TypeMux
	engine         consensus.Engine
	accountManager *accounts.Manager

	bloomRequests chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	bloomIndexer  *core.ChainIndexer             // Bloom indexer operating during block imports

	APIBackend *TrueAPIBackend

	//miner     *miner.Miner
	gasPrice *big.Int

	networkID     uint64
	netRPCService *taiapi.PublicNetAPI

	pbftServer *tbft.Node

	lock sync.RWMutex // Protects the variadic fields (e.g. gas price and etherbase)
}

func (s *Taiyuechain) AddLesServer(ls LesServer) {
	s.lesServer = ls
	ls.SetBloomBitsIndexer(s.bloomIndexer)
}

// New creates a new Taiyuechain object (including the
// initialisation of the common Taiyuechain object)
func New(ctx *node.ServiceContext, config *Config, p2pCert []byte) (*Taiyuechain, error) {
	if config.SyncMode == downloader.LightSync {
		return nil, errors.New("can't run yue.Taiyuechain in light sync mode, use les.LightTaiYueChain")
	}

	if !config.SyncMode.IsValid() {
		return nil, fmt.Errorf("invalid sync mode %d", config.SyncMode)
	}
	chainDb, err := CreateDB(ctx, config, "chaindata")
	//chainDb, err := CreateDB(ctx, config, path)
	if err != nil {
		return nil, err
	}

	chainConfig, _, genesisErr := core.SetupGenesisBlock(chainDb, config.Genesis)
	if _, ok := genesisErr.(*params.ConfigCompatError); genesisErr != nil && !ok {
		return nil, genesisErr
	}

	log.Info("Initialised chain configuration", "config", chainConfig)

	NewCIMList := cim.NewCIMList(uint8(crypto.CryptoType))



	yue := &Taiyuechain{
		config:         config,
		chainDb:        chainDb,
		chainConfig:    chainConfig,
		eventMux:       ctx.EventMux,
		accountManager: ctx.AccountManager,
		engine:         CreateConsensusEngine(ctx, &ethash.Config{PowMode: ethash.ToMinervaMode(config.MinervaMode)},NewCIMList),
		shutdownChan:   make(chan bool),
		networkID:      config.NetworkId,
		gasPrice:       config.GasPrice,
		bloomRequests:  make(chan chan *bloombits.Retrieval),
		bloomIndexer:   NewBloomIndexer(chainDb, params.BloomBitsBlocks),
	}

	log.Info("Initialising Taiyuechain protocol", "versions", ProtocolVersions, "network", config.NetworkId)

	/*if !config.SkipBcVersionCheck {
		bcVersion := rawdb.ReadDatabaseVersion(chainDb)
		if bcVersion != core.BlockChainVersion && bcVersion != 0 {
			return nil, fmt.Errorf("Blockchain DB version mismatch (%d / %d). Run taiyue upgradedb.\n", bcVersion, core.BlockChainVersion)
		}
		rawdb.WriteDatabaseVersion(chainDb, core.BlockChainVersion)
	}*/
	var (
		vmConfig    = vm.Config{EnablePreimageRecording: config.EnablePreimageRecording}
		cacheConfig = &core.CacheConfig{Deleted: config.DeletedState, Disabled: config.NoPruning, TrieNodeLimit: config.TrieCache, TrieTimeLimit: config.TrieTimeout}
	)
	//NewCIMList := cim.NewCIMList(yue.config.CryptoType)

	yue.blockchain, err = core.NewBlockChain(chainDb, cacheConfig, yue.chainConfig, yue.engine, vmConfig, NewCIMList)
	if err != nil {
		return nil, err
	}

	gensysExra :=yue.blockchain.Genesis().Extra()
	pmSt := false
	pmcC := false
	if uint8(gensysExra[3]) >0 {
		pmSt = true
	}
	if uint8(gensysExra[4]) > 0{
		pmcC = true
	}
	vm.SetPermConfig(pmSt,pmcC)


	//init cert list to
	// need init cert list to statedb
	stateDB, err := yue.blockchain.State()
	if err != nil {
		return nil, err
	}

	err = NewCIMList.InitCertAndPermission(yue.blockchain.CurrentBlock().Number(), stateDB)
	if err != nil {
		panic(err)
	}

	// Rewind the chain in case of an incompatible config upgrade.
	/*if compat, ok := genesisErr.(*params.ConfigCompatError); ok {
		log.Warn("Rewinding chain to upgrade configuration", "err", compat)
		yue.blockchain.SetHead(compat.RewindTo)
		rawdb.WriteChainConfig(chainDb, genesisHash, chainConfig)
	}*/

	//  rewind snail if case of incompatible config
	/*if compat, ok := genesisErr.(*params.ConfigCompatError); ok {
		log.Warn("Rewinding snail chain to upgrade configuration", "err", compat)
		yue.snailblockchain.SetHead(compat.RewindTo)
		rawdb.WriteChainConfig(chainDb, genesisHash, chainConfig)
	}*/

	yue.bloomIndexer.Start(yue.blockchain)

	if config.TxPool.Journal != "" {
		config.TxPool.Journal = ctx.ResolvePath(config.TxPool.Journal)
	}
	config.TxPool.CimList = NewCIMList

	yue.txPool = core.NewTxPool(config.TxPool, yue.chainConfig, yue.blockchain)

	yue.election = elect.NewElection(yue.blockchain, yue.config)

	yue.engine.SetElection(yue.election)

	//coinbase, _ := yue.Etherbase()

	cacheLimit := cacheConfig.TrieCleanLimit //+ cacheConfig.TrieDirtyLimit
	checkpoint := config.Checkpoint
	//TODO neo
	/*if checkpoint == nil {
		checkpoint = params.TrustedCheckpoints[genesisHash]
	}*/
	yue.agent = NewPbftAgent(yue, yue.chainConfig, yue.engine, yue.election,
		NewCIMList, config.MinerGasFloor, config.MinerGasCeil)
	if yue.protocolManager, err = NewProtocolManager(yue.chainConfig, checkpoint, config.SyncMode, config.NetworkId, yue.eventMux, yue.txPool, yue.engine, yue.blockchain, chainDb, yue.agent, cacheLimit, config.Whitelist, NewCIMList, p2pCert); err != nil {
		return nil, err
	}

	//committeeKey, err := crypto.ToECDSA(yue.config.CommitteeKey)
	//if err == nil {
	//	yue.miner.SetElection(yue.config.EnableElection, crypto.FromECDSAPub(&committeeKey.PublicKey))
	//}

	yue.APIBackend = &TrueAPIBackend{yue, nil}
	gpoParams := config.GPO
	if gpoParams.Default == nil {
		gpoParams.Default = config.GasPrice
	}
	yue.APIBackend.gpo = gasprice.NewOracle(yue.APIBackend, gpoParams)
	return yue, nil
}

func makeExtraData(extra []byte) []byte {
	if len(extra) == 0 {
		// create default extradata
		extra, _ = rlp.EncodeToBytes([]interface{}{
			uint(params.VersionMajor<<16 | params.VersionMinor<<8 | params.VersionPatch),
			"taiyue",
			runtime.Version(),
			runtime.GOOS,
		})
	}
	if uint64(len(extra)) > params.MaximumExtraDataSize {
		log.Warn("Miner extra data exceed limit", "extra", hexutil.Bytes(extra), "limit", params.MaximumExtraDataSize)
		extra = nil
	}
	return extra
}

// CreateDB creates the chain database.
func CreateDB(ctx *node.ServiceContext, config *Config, name string) (yuedb.Database, error) {
	db, err := ctx.OpenDatabase(name, config.DatabaseCache, config.DatabaseHandles)
	if err != nil {
		return nil, err
	}
	if db, ok := db.(*yuedb.LDBDatabase); ok {
		db.Meter("yue/db/chaindata/")
	}
	return db, nil
}

// CreateConsensusEngine creates the required type of consensus engine instance for an Taiyuechain service
func CreateConsensusEngine(ctx *node.ServiceContext, config *ethash.Config,cimList *cim.CimList) consensus.Engine {
	// Otherwise assume proof-of-work
	switch config.PowMode {
	case ethash.ModeFake:
		log.Info("-----Fake mode")
		log.Warn("Ethash used in fake mode")
		return ethash.NewFaker(cimList)
	case ethash.ModeTest:
		log.Warn("Ethash used in test mode")
		return ethash.NewTester(cimList)
	case ethash.ModeShared:
		log.Warn("Ethash used in shared mode")
		return ethash.NewShared(cimList)
	default:
		engine := ethash.New(ethash.Config{
			PowMode: config.PowMode,
		},cimList)
		//engine.SetThreads(-1) // Disable CPU mining
		return engine
	}
}

// APIs return the collection of RPC services the yue package offers.
// NOTE, some of these services probably need to be moved to somewhere else.
func (s *Taiyuechain) APIs() []rpc.API {
	apis := taiapi.GetAPIs(s.APIBackend)

	// Append any APIs exposed explicitly by the consensus engine
	apis = append(apis, s.engine.APIs(s.BlockChain())...)

	// Append yue	APIs and  Eth APIs
	namespaces := []string{"etrue", "yue"}
	for _, name := range namespaces {
		apis = append(apis, []rpc.API{
			{
				Namespace: name,
				Version:   "1.0",
				Service:   NewPublicTaiyueChainAPI(s),
				Public:    true,
			}, {
				Namespace: name,
				Version:   "1.0",
				Service:   downloader.NewPublicDownloaderAPI(s.protocolManager.downloader, s.eventMux),
				Public:    true,
			}, {
				Namespace: name,
				Version:   "1.0",
				Service:   filters.NewPublicFilterAPI(s.APIBackend, false),
				Public:    true,
			},
		}...)
	}
	// Append all the local APIs and return
	return append(apis, []rpc.API{
		{
			Namespace: "admin",
			Version:   "1.0",
			Service:   NewPrivateAdminAPI(s),
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPublicDebugAPI(s),
			Public:    true,
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPrivateDebugAPI(s.chainConfig, s),
		}, {
			Namespace: "net",
			Version:   "1.0",
			Service:   s.netRPCService,
			Public:    true,
		}, {
			Namespace: "cpm",
			Version:   "1.0",
			Service:   taiapi.NewPublicCertAPI(s.APIBackend),
			Public:    true,
		},
	}...)
}

func (s *Taiyuechain) ResetWithGenesisBlock(gb *types.Block) {
	s.blockchain.ResetWithGenesisBlock(gb)
}

func (s *Taiyuechain) ResetWithFastGenesisBlock(gb *types.Block) {
	s.blockchain.ResetWithGenesisBlock(gb)
}

func (s *Taiyuechain) Etherbase() common.Address {
	etherbase := s.agent.committeeNode.Coinbase
	return etherbase
}

func (s *Taiyuechain) PbftAgent() *PbftAgent              { return s.agent }
func (s *Taiyuechain) AccountManager() *accounts.Manager  { return s.accountManager }
func (s *Taiyuechain) BlockChain() *core.BlockChain       { return s.blockchain }
func (s *Taiyuechain) Config() *Config                    { return s.config }
func (s *Taiyuechain) TxPool() *core.TxPool               { return s.txPool }
func (s *Taiyuechain) EventMux() *event.TypeMux           { return s.eventMux }
func (s *Taiyuechain) Engine() consensus.Engine           { return s.engine }
func (s *Taiyuechain) ChainDb() yuedb.Database            { return s.chainDb }
func (s *Taiyuechain) IsListening() bool                  { return true } // Always listening
func (s *Taiyuechain) EthVersion() int                    { return int(s.protocolManager.SubProtocols[0].Version) }
func (s *Taiyuechain) NetVersion() uint64                 { return s.networkID }
func (s *Taiyuechain) Downloader() *downloader.Downloader { return s.protocolManager.downloader }

// Protocols implements node.Service, returning all the currently configured
// network protocols to start.
func (s *Taiyuechain) Protocols() []p2p.Protocol {

	if s.lesServer == nil {
		return s.protocolManager.SubProtocols
	}
	return append(s.protocolManager.SubProtocols, s.lesServer.Protocols()...)
	/*protos := make([]p2p.Protocol, len(ProtocolVersions))
	for i, vsn := range ProtocolVersions {
		protos[i] = s.protocolManager.makeProtocol(vsn)
		//protos[i].Attributes = []enr.Entry{s.currentEthEntry()}
		//protos[i].DialCandidates = s.dialCandiates
	}
	if s.lesServer != nil {
		protos = append(protos, s.lesServer.Protocols()...)
	}
	return protos*/
}

// Start implements node.Service, starting all internal goroutines needed by the
// Taiyuechain protocol implementation.
func (s *Taiyuechain) Start(srvr *p2p.Server) error {

	// Start the bloom bits servicing goroutines
	s.startBloomHandlers()

	// Start the RPC service
	s.netRPCService = taiapi.NewPublicNetAPI(srvr, s.NetVersion())

	// Figure out a max peers count based on the server limits
	maxPeers := srvr.MaxPeers
	// Start the networking layer and the light server if requested
	s.protocolManager.Start(maxPeers)
	if s.lesServer != nil {
		s.lesServer.Start(srvr)
	}
	s.startPbftServer()
	if s.pbftServer == nil {
		log.Error("start pbft server failed.")
		return errors.New("start pbft server failed.")
	}
	s.agent.server = s.pbftServer
	log.Info("Start", "server", s.agent.server, "SyncMode", s.config.SyncMode)
	s.agent.Start()

	s.election.Start()

	// Start the networking layer and the light server if requested
	/*s.protocolManager.Start2(maxPeers)*/

	return nil
}

// Stop implements node.Service, terminating all internal goroutines used by the
// Taiyuechain protocol.
func (s *Taiyuechain) Stop() error {
	s.stopPbftServer()
	s.bloomIndexer.Close()
	s.blockchain.Stop()
	s.protocolManager.Stop()
	if s.lesServer != nil {
		s.lesServer.Stop()
	}
	s.txPool.Stop()
	s.eventMux.Stop()

	s.chainDb.Close()
	close(s.shutdownChan)

	return nil
}

func (s *Taiyuechain) startPbftServer() error {

	priv, err := crypto.ToECDSA(s.config.CommitteeKey)
	if err != nil {
		return err
	}

	cfg := config.DefaultConfig()
	cfg.P2P.ListenAddress1 = "tcp://0.0.0.0:" + strconv.Itoa(s.config.Port)
	cfg.P2P.ListenAddress2 = "tcp://0.0.0.0:" + strconv.Itoa(s.config.StandbyPort)

	n1, err := tbft.NewNode(cfg, "1", priv, s.agent)
	if err != nil {
		return err
	}
	s.pbftServer = n1
	return n1.Start()
}

func (s *Taiyuechain) stopPbftServer() error {
	if s.pbftServer != nil {
		s.pbftServer.Stop()
	}
	return nil
}
