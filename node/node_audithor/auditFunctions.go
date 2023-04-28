package node_audithor

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/networking"
	"bavovnacoin/node/node_controller/node_settings"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type ReorgLogs struct {
	Time   int64
	Height uint64
}

var reorgData []ReorgLogs

func GetReorgData() {
	jsonFile, err := os.Open("data / AudNode / reorgData.json")
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &reorgData)
	jsonFile.Close()
}

func SaveReorgData() {
	byteData, _ := json.MarshalIndent(reorgData, "", "    ")
	os.WriteFile("data/AudNode/reorgData.json", byteData, 0777)
}

func RecieveBlocks() ([]blockchain.Block, int) {
	var conn networking.Connection
	isConn, _ := conn.EstablishAddresses(node_settings.Settings.OtherNodesAddresses[0], -1, "")
	if !isConn {
		return nil, 1
	}

	blocksOnH, height, res := conn.RequestBlocks(blockchain.BcLength)
	if !res || len(blocksOnH) == 0 || height == blockchain.BcLength {
		return nil, 2
	}

	var blocks []blockchain.Block
	for _, bl := range blocksOnH[0].Blocks {
		blocks = append(blocks, bl.Block)
	}

	return blocks, 0
}

func ReorgTests(blocks []blockchain.Block) {
	if len(blocks) > 1 {
		println("Reorganization happened")
		reorgData = append(reorgData, ReorgLogs{Time: time.Now().UTC().Unix(), Height: blockchain.BcLength - uint64(len(blocks))})
		SaveReorgData()
	}

	for i := 0; i < len(blocks); i++ {
		blockchain.AddBlockToBlockchain(blocks[i], 0, true)
		blockchain.IncrBcHeight(0)
		log.Println("Block is added to blockchain. Current height: " + fmt.Sprint(int(blockchain.BcLength)) + "\n")
	}
}
