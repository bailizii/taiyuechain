// Code generated by github.com/fjl/gencodec. DO NOT EDIT.

package yue

import (
	"math/big"
	"time"

	"github.com/taiyuechain/taiyuechain/common"
	"github.com/taiyuechain/taiyuechain/common/hexutil"
	"github.com/taiyuechain/taiyuechain/core"
	"github.com/taiyuechain/taiyuechain/params"
	"github.com/taiyuechain/taiyuechain/yue/downloader"
	"github.com/taiyuechain/taiyuechain/yue/gasprice"
)

// MarshalTOML marshals as TOML.
func (c Config) MarshalTOML() (interface{}, error) {
	type Config struct {
		Genesis                 *core.Genesis `toml:",omitempty"`
		NetworkId               uint64
		SyncMode                downloader.SyncMode
		NoPruning               bool
		Whitelist               map[uint64]common.Hash `toml:"-"`
		SkipBcVersionCheck      bool                   `toml:"-"`
		DatabaseHandles         int                    `toml:"-"`
		DatabaseCache           int
		TrieCache               int
		TrieTimeout             time.Duration
		MinervaMode             int
		Host                    string
		CommitteeKey            hexutil.Bytes
		CommitteeBase           common.Address
		NodeCert                hexutil.Bytes
		Port                    int
		StandbyPort             int
		NodeType                bool
		GasPrice                *big.Int `toml:",omitempty"`
		MinerGasCeil            uint64
		MinerGasFloor           uint64
		TxPool                  core.TxPoolConfig
		GPO                     gasprice.Config
		EnablePreimageRecording bool
		EnableElection          bool
		DocRoot                 string                    `toml:"-"`
		Checkpoint              *params.TrustedCheckpoint `toml:",omitempty"`
	}
	var enc Config
	enc.Genesis = c.Genesis
	enc.NetworkId = c.NetworkId
	enc.SyncMode = c.SyncMode
	enc.NoPruning = c.NoPruning
	enc.Whitelist = c.Whitelist
	enc.SkipBcVersionCheck = c.SkipBcVersionCheck
	enc.DatabaseHandles = c.DatabaseHandles
	enc.DatabaseCache = c.DatabaseCache
	enc.TrieCache = c.TrieCache
	enc.MinervaMode = c.MinervaMode
	enc.TrieTimeout = c.TrieTimeout
	enc.Host = c.Host
	enc.Port = c.Port
	enc.MinerGasCeil = c.MinerGasCeil
	enc.MinerGasFloor = c.MinerGasFloor
	enc.StandbyPort = c.StandbyPort
	enc.CommitteeKey = c.CommitteeKey
	enc.CommitteeBase = c.CommitteeBase
	enc.NodeCert = c.NodeCert
	enc.NodeType = c.NodeType
	enc.GasPrice = c.GasPrice
	enc.TxPool = c.TxPool
	enc.GPO = c.GPO
	enc.EnablePreimageRecording = c.EnablePreimageRecording
	enc.EnableElection = c.EnableElection
	enc.DocRoot = c.DocRoot
	enc.Checkpoint = c.Checkpoint
	return &enc, nil
}

// UnmarshalTOML unmarshals from TOML.
func (c *Config) UnmarshalTOML(unmarshal func(interface{}) error) error {
	type Config struct {
		Genesis                 *core.Genesis `toml:",omitempty"`
		NetworkId               *uint64
		SyncMode                *downloader.SyncMode
		NoPruning               *bool
		Whitelist               map[uint64]common.Hash `toml:"-"`
		SkipBcVersionCheck      *bool                  `toml:"-"`
		DatabaseHandles         *int                   `toml:"-"`
		DatabaseCache           *int
		TrieCache               *int
		MinervaMode             *int
		Host                    *string
		Port                    *int
		StandbyPort             *int
		MinerGasCeil            *uint64
		MinerGasFloor           *uint64
		CommitteeKey            *hexutil.Bytes
		CommitteeBase           *common.Address
		NodeCert                *hexutil.Bytes
		TrieTimeout             *time.Duration
		NodeType                *bool
		TxPool                  *core.TxPoolConfig
		GasPrice                *big.Int `toml:",omitempty"`
		GPO                     *gasprice.Config
		EnablePreimageRecording *bool
		EnableElection          *bool
		DocRoot                 *string                   `toml:"-"`
		Checkpoint              *params.TrustedCheckpoint `toml:",omitempty"`
	}
	var dec Config
	if err := unmarshal(&dec); err != nil {
		return err
	}
	if dec.Genesis != nil {
		c.Genesis = dec.Genesis
	}
	if dec.NetworkId != nil {
		c.NetworkId = *dec.NetworkId
	}
	if dec.SyncMode != nil {
		c.SyncMode = *dec.SyncMode
	}
	if dec.NoPruning != nil {
		c.NoPruning = *dec.NoPruning
	}
	if dec.Whitelist != nil {
		c.Whitelist = dec.Whitelist
	}
	if dec.SkipBcVersionCheck != nil {
		c.SkipBcVersionCheck = *dec.SkipBcVersionCheck
	}
	if dec.DatabaseHandles != nil {
		c.DatabaseHandles = *dec.DatabaseHandles
	}
	if dec.DatabaseCache != nil {
		c.DatabaseCache = *dec.DatabaseCache
	}
	if dec.TrieCache != nil {
		c.TrieCache = *dec.TrieCache
	}
	if dec.MinervaMode != nil {
		c.MinervaMode = *dec.MinervaMode
	}
	if dec.TrieTimeout != nil {
		c.TrieTimeout = *dec.TrieTimeout
	}
	if dec.Host != nil {
		c.Host = *dec.Host
	}
	if dec.Port != nil {
		c.Port = *dec.Port
	}
	if dec.StandbyPort != nil {
		c.StandbyPort = *dec.StandbyPort
	}
	if dec.CommitteeKey != nil {
		c.CommitteeKey = *dec.CommitteeKey
	}
	if dec.CommitteeBase != nil {
		c.CommitteeBase = *dec.CommitteeBase
	}
	if dec.NodeCert != nil {
		c.NodeCert = *dec.NodeCert
	}
	if dec.NodeType != nil {
		c.NodeType = *dec.NodeType
	}
	if dec.TxPool != nil {
		c.TxPool = *dec.TxPool
	}
	if dec.GasPrice != nil {
		c.GasPrice = dec.GasPrice
	}
	if dec.MinerGasCeil != nil {
		c.MinerGasCeil = *dec.MinerGasCeil
	}
	if dec.MinerGasFloor != nil {
		c.MinerGasFloor = *dec.MinerGasFloor
	}
	if dec.GPO != nil {
		c.GPO = *dec.GPO
	}
	if dec.EnablePreimageRecording != nil {
		c.EnablePreimageRecording = *dec.EnablePreimageRecording
	}
	if dec.EnableElection != nil {
		c.EnableElection = *dec.EnableElection
	}
	if dec.DocRoot != nil {
		c.DocRoot = *dec.DocRoot
	}
	if dec.Checkpoint != nil {
		c.Checkpoint = dec.Checkpoint
	}
	return nil
}