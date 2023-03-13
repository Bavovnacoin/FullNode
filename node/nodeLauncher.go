package node

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/dbController"
	"bavovnacoin/networking"
	"bavovnacoin/node_controller"
	"bavovnacoin/node_controller/command_executor"
	"bavovnacoin/txo"
	"fmt"
)

func StartRPC() {
	isRpcStarted, err := networking.StartRPCListener()
	if !isRpcStarted {
		fmt.Println("Can't start RPC listener")
		fmt.Println(err)
	} else {
		fmt.Println("RPC listener started")
	}
}

func NodeProcess() {
	blockchain.RestoreMempool()
	txo.RestoreCoinDatabase()
	blockchain.InitBlockchain()

	for command_executor.Node_working {
		AddBlock(true)
	}
}

func Launch() {
	dbController.DB.OpenDb()
	StartRPC()
	go NodeProcess()
	node_controller.CommandHandler()
	blockchain.BackTransactionsToMempool()
	blockchain.WriteMempoolData()
	dbController.DB.CloseDb()
}
