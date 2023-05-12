package node_settings

import (
	"bavovnacoin/byteArr"
	"bavovnacoin/cryption"
	"bavovnacoin/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"runtime"
)

var Settings NodeSettings
var settingsFileName string = "node_settings.json"

type NodeSettings struct {
	TxMinFee            uint64     // Lowest fee of tx to be verified
	MiningThreads       uint       // Threads for mining count
	NodeType            uint       // For now there's two types: full node, audithor
	OtherNodesAddresses [][]string // Addresses of other nodes to communicate
	MyAddress           string
	PrivKey             []byte
	PrivKeyDecrypted    []byte `json:"-"`
	HashPass            byteArr.ByteArr
	RewardAddress       string
	RPCip               string

	NodeTypesNames []string `json:"-"`
}

func (ns *NodeSettings) InitSettingsValues() {
	ns.NodeTypesNames = []string{"Validator node", "Audithor node"}
}

func (ns *NodeSettings) GetSettings() bool {
	jsonFile, err := os.Open(settingsFileName)
	defer jsonFile.Close()

	if err != nil {
		*ns = NodeSettings{TxMinFee: 3, MiningThreads: 0, NodeType: 0}
		ns.WriteSettings()
		return false
	} else {
		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &ns)
	}
	return true
}

func (ns *NodeSettings) DecryptPrivKey(password string) {
	ns.PrivKeyDecrypted = []byte(cryption.AES_decrypt(string(ns.PrivKey), password))
}

func (ns *NodeSettings) WriteSettings() {
	byteData, _ := json.MarshalIndent(ns, "", "    ")
	os.WriteFile(settingsFileName, byteData, 0777)
}

func (ns *NodeSettings) GetThreadsAmmountForMining() uint {
	if ns.MiningThreads == 0 {
		if uint(runtime.NumCPU()-4) <= 0 {
			return 1
		}
		return uint(runtime.NumCPU() - 4)
	}
	return ns.MiningThreads
}

func (ns *NodeSettings) GetMaxThreadsAmmount() uint64 {
	thrAmmount := uint64(runtime.NumCPU() - 4)
	if thrAmmount <= 0 {
		return 1
	}
	return thrAmmount
}

func (ns *NodeSettings) SetMiningThreads(threads uint) bool {
	if threads < 0 || threads > uint(ns.GetMaxThreadsAmmount()) {
		return false
	}
	ns.MiningThreads = threads
	return true
}

func (ns *NodeSettings) ThreadsForMiningToString() string {
	if ns.MiningThreads == 0 {
		return fmt.Sprintf("MAX (%d)", ns.GetThreadsAmmountForMining())
	}
	return fmt.Sprintf("%d", ns.MiningThreads)
}

func (ns *NodeSettings) IsAddressAdded(address string) bool {
	for i := 0; i < len(ns.OtherNodesAddresses); i++ {
		if ns.OtherNodesAddresses[i][0] == address {
			return true
		}
	}
	return false
}

func (ns *NodeSettings) AddAddress(address []string) bool {
	if !ns.IsAddressAdded(address[0]) {
		ns.OtherNodesAddresses = append(ns.OtherNodesAddresses, address)
		return true
	}
	return false
}

func (ns *NodeSettings) RemAddress(ind int) {
	ns.OtherNodesAddresses = append(ns.OtherNodesAddresses[:ind], ns.OtherNodesAddresses[ind+1:]...)
}

func (ns *NodeSettings) GetRewAddress() string {
	if ns.RewardAddress == "" {
		return "none"
	}
	return ns.RewardAddress
}

func (ns *NodeSettings) IsRewAddrWalid(rewAddr string) bool {
	_, res := new(big.Int).SetString(rewAddr, 16)
	if len(rewAddr) != 40 || !res {
		return false
	}
	return true
}

func (ns *NodeSettings) SetPrivKey(password string) {
	ecdsa.InitValues()
	ns.PrivKey = []byte(cryption.AES_encrypt(ecdsa.GenPrivKey(), password))
	ns.WriteSettings()
}

func (ns *NodeSettings) GetPrivKey() []byte {
	return Settings.PrivKeyDecrypted
}

func (ns *NodeSettings) GetRpcAddr() string {
	return Settings.RPCip
}
