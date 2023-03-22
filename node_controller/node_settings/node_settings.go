package node_settings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
)

var Settings NodeSettings
var settingsFileName string = "node_settings.json"

type NodeSettings struct {
	TxMinFee            uint64   // Lowest fee of tx to be verified
	MiningThreads       uint     // Threads for mining count
	NodeType            uint     // For now there's two types: full node, audithor
	OtherNodesAddresses []string // Addresses of other nodes to communicate

	NodeTypesNames []string `json:"-"`
}

func (ns *NodeSettings) InitSettingsValues() {
	ns.NodeTypesNames = []string{"Full node", "Audithor node"}
}

func (ns *NodeSettings) GetSettings() {
	jsonFile, err := os.Open(settingsFileName)
	defer jsonFile.Close()

	if err != nil {
		*ns = NodeSettings{TxMinFee: 3, MiningThreads: 0, NodeType: 0, OtherNodesAddresses: []string{"188.163.8.146:25565"}}
		ns.WriteSettings()
	} else {
		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &ns)
	}
}

func (ns *NodeSettings) WriteSettings() {
	byteData, _ := json.MarshalIndent(ns, "", "    ")
	os.WriteFile(settingsFileName, byteData, 0777)
}

func getThreadsAmmount() int {
	thrAmmount := runtime.NumCPU() - 4
	if thrAmmount <= 0 {
		return 1
	}
	return thrAmmount
}

func (ns *NodeSettings) SetMiningThreads(threads uint) bool {
	if threads < 0 || threads > uint(getThreadsAmmount()) {
		return false
	}
	ns.MiningThreads = threads
	return true
}

func (ns *NodeSettings) ThreadsForMiningToString() string {
	if ns.MiningThreads == 0 {
		return fmt.Sprintf("MAX (%d)", getThreadsAmmount())
	}
	return fmt.Sprintf("%d", ns.MiningThreads)
}

func (ns *NodeSettings) IsAddressAdded(address string) bool {
	for i := 0; i < len(ns.OtherNodesAddresses); i++ {
		if ns.OtherNodesAddresses[i] == address {
			return true
		}
	}
	return false
}

func (ns *NodeSettings) IsAddressValid(address string) bool {
	isAddrMatch, _ := regexp.MatchString("^(?:http(s)?:\\/\\/)?[\\w.-]+(?:\\.[\\w\\.-]+)+[\\w\\-\\._~:/?#[\\]@!\\$&'\\(\\)\\*\\+,;=.]+$", address)
	return isAddrMatch
}

func (ns *NodeSettings) AddAddress(address string) bool {
	if ns.IsAddressValid(address) && !ns.IsAddressAdded(address) {
		ns.OtherNodesAddresses = append(ns.OtherNodesAddresses, address)
		return true
	}
	return false
}

func (ns *NodeSettings) RemAddress(ind int) {
	ns.OtherNodesAddresses = append(ns.OtherNodesAddresses[:ind], ns.OtherNodesAddresses[ind+1:]...)
}
