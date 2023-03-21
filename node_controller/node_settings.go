package node_controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var settingsFileName string = "node_settings.json"

type NodeSettings struct {
	TxAshFee            uint64   // Lowest fee of tx to be verified
	MiningThreads       uint     // Threads for mining count
	NodeType            uint     // For now there's two types: full node, audithor
	OtherNodesAddresses []string // Addresses of other nodes to communicate
}

func (ns *NodeSettings) GetSettings() {
	jsonFile, err := os.Open(settingsFileName)
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &ns)
	jsonFile.Close()
}

func (ns *NodeSettings) WriteSettings() {
	byteData, _ := json.MarshalIndent(ns, "", "    ")
	os.WriteFile(settingsFileName, byteData, 0777)
}
