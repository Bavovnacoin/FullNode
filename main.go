package main

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/dbController"
	"bavovnacoin/hashing"
)

func GetBlockHashes(args ...uint64) {
	dbController.DB.OpenDb()
	defer dbController.DB.CloseDb()
	blockchain.InitBlockchain()

	for _, ind := range args {
		b, _ := blockchain.GetBlock(ind)
		println(hashing.SHA1(blockchain.BlockToString(b)), blockchain.ValidateBlock(b, int(ind), true, true))
	}
}

func main() {
	// node.Launch()

	GetBlockHashes(2, 5)
}
