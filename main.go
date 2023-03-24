package main

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/dbController"
	"bavovnacoin/hashing"
)

func GetBlockHashes(args ...uint64) {
	dbController.DB.OpenDb()
	defer dbController.DB.CloseDb()

	for _, ind := range args {
		b, _ := blockchain.GetBlock(ind)
		println(hashing.SHA1(blockchain.BlockToString(b)), blockchain.ValidateBlock(b, int(ind), true, true))
	}
}

func main() {
	// node.Launch()

	GetBlockHashes(0, 5)
}
