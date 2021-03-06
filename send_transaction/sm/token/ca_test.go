package test

import (
	"encoding/hex"
	"fmt"
	"github.com/taiyuechain/taiyuechain/cim"
	"github.com/taiyuechain/taiyuechain/core/vm"
	"github.com/taiyuechain/taiyuechain/crypto"
	"math/big"
	"os"
	"testing"

	"github.com/taiyuechain/taiyuechain/core"
	"github.com/taiyuechain/taiyuechain/core/state"
	"github.com/taiyuechain/taiyuechain/core/types"
	"github.com/taiyuechain/taiyuechain/log"
)

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlInfo, log.StreamHandler(os.Stderr, log.TerminalFormat(false))))
}

//neo test cacert contract
func TestAllCaCert(t *testing.T) {
	// Create a helper to check if a gas allowance results in an executable transaction
	executable := func(number uint64, gen *core.BlockGen, fastChain *core.BlockChain, header *types.Header, statedb *state.StateDB, cimList *cim.CimList) {
		sendTranction(number, gen, statedb, pAccount1, pAccount2, big.NewInt(6000000000000000000), pKey1, signer, nil, header, p2p1Byte)
		cert44 := pbft5Byte
		if number == 30 {
			fmt.Println("amount",getCaAmount(statedb,header.Number))
		}

		sendMultiProposalTranscation(number, gen, pAccount2, cert44, pbft1Byte,pub1, true, pKey2, signer, statedb, fastChain, abiCA, nil, p2p2Byte)
		sendMultiProposalTranscation(number, gen, pAccount2, cert44, pbft2Byte,pub2, true, pKey2, signer, statedb, fastChain, abiCA, nil, p2p2Byte)

		if number == 1050 {
			fmt.Println("amount",getCaAmount(statedb,header.Number),"number",header.Number)
			//if err := cimList.VerifyRootCert(cert44); err != nil {
			//	fmt.Println("TestAllCaCert err", err," header ",header.Number)
			//}
		}

		sendMultiProposalTranscation(number-26-1000, gen, pAccount2, cert44, pbft1Byte,pub1, false, pKey2, signer, statedb, fastChain, abiCA, nil, p2p2Byte)
		sendMultiProposalTranscation(number-27-1000, gen, pAccount2, cert44, pbft2Byte,pub2, false, pKey2, signer, statedb, fastChain, abiCA, nil, p2p2Byte)
		sendMultiProposalTranscation(number-28-1000, gen, pAccount2, cert44, pbft3Byte,pub1, false, pKey2, signer, statedb, fastChain, abiCA, nil, p2p2Byte)
		if number == 2050 {
			fmt.Println("amount",getCaAmount(statedb,header.Number),"number",header.Number)
		}
	}
	newTestPOSManager(50, executable)
	fmt.Println("staking addr", types.CACertListAddress)
}

func getCaAmount(state *state.StateDB, number *big.Int) uint64  {
	caCertList := vm.NewCACertList()
	err := caCertList.LoadCACertList(state, types.CACertListAddress)

	if err != nil {
		log.Error("Staking load error", "error", err)
	}

	return caCertList.GetCaCertAmount(types.GetEpochIDFromHeight(number).Uint64())
}

func TestGetAddress(t *testing.T) {
	// Create a helper to check if a gas allowance results in an executable transaction
	skey, _ := crypto.HexToECDSA("7631a11e9d28563cdbcf96d581e4b9a19e53ad433a53c25a9f18c74ddf492f75")
	saddr := crypto.PubkeyToAddress(skey.PublicKey)
	skey2, _ := crypto.HexToECDSA("bab8dbdcb4d974eba380ff8b2e459efdb6f8240e5362e40378de3f9f5f1e67bb")
	saddr2 := crypto.PubkeyToAddress(skey2.PublicKey)
	//103
	skey3, _ := crypto.HexToECDSA("122d186b77a030e04f5654e13d934b21af2aac03b942c3ecda4632364d81cbab")
	saddr3 := crypto.PubkeyToAddress(skey3.PublicKey)
	//104
	skey4, _ := crypto.HexToECDSA("fe44cbc0e164092a6746bd57957422ab165c009d0299c7639a2f4d290317f20f")
	saddr4 := crypto.PubkeyToAddress(skey4.PublicKey)

	fmt.Println("saddr", crypto.AddressToHex(saddr), "saddr2", crypto.AddressToHex(saddr2), "saddr3", crypto.AddressToHex(saddr3), "saddr4 ", crypto.AddressToHex(saddr4))

	pub101 := hex.EncodeToString(crypto.FromECDSAPub(&skey.PublicKey))
	pub102 := hex.EncodeToString(crypto.FromECDSAPub(&skey2.PublicKey))
	pub103 := hex.EncodeToString(crypto.FromECDSAPub(&skey3.PublicKey))
	pub104 := hex.EncodeToString(crypto.FromECDSAPub(&skey4.PublicKey))
	fmt.Println("pub101", pub101, "pub102", pub102, "pub103", pub103, "pub104 ", pub104)

	// Create a helper to check if a gas allowance results in an executable transaction
	skey, _ = crypto.HexToECDSA("d5939c73167cd3a815530fd8b4b13f1f5492c1c75e4eafb5c07e8fb7f4b09c7c")
	saddr = crypto.PubkeyToAddress(skey.PublicKey)
	skey2, _ = crypto.HexToECDSA("ea4297749d514cc476fe971a7fe20100cbd29f010864341b3e624e8744d46cec")
	saddr2 = crypto.PubkeyToAddress(skey2.PublicKey)
	//103
	skey3, _ = crypto.HexToECDSA("86937006ac1e6e2c846e160d93f86c0d63b0fcefc39a46e9eaeb65188909fbdc")
	saddr3 = crypto.PubkeyToAddress(skey3.PublicKey)
	//104
	skey4, _ = crypto.HexToECDSA("cbddcbecd252a8586a4fd759babb0cc77f119d55f38bc7f80a708e75964dd801")
	saddr4 = crypto.PubkeyToAddress(skey4.PublicKey)

	fmt.Println("saddr", crypto.AddressToHex(saddr), "saddr2", crypto.AddressToHex(saddr2), "saddr3", crypto.AddressToHex(saddr3), "saddr4 ", crypto.AddressToHex(saddr4))

	pub101 = hex.EncodeToString(crypto.FromECDSAPub(&skey.PublicKey))
	pub102 = hex.EncodeToString(crypto.FromECDSAPub(&skey2.PublicKey))
	pub103 = hex.EncodeToString(crypto.FromECDSAPub(&skey3.PublicKey))
	pub104 = hex.EncodeToString(crypto.FromECDSAPub(&skey4.PublicKey))
	fmt.Println("pub101", pub101, "pub102", pub102, "pub103", pub103, "pub104 ", pub104)
}
