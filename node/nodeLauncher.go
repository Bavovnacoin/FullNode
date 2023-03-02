package node

import (
	"bavovnacoin/account"
	"bavovnacoin/blockchain"
	"bavovnacoin/dbController"
	"bavovnacoin/networking"
	"bavovnacoin/node_controller"
	"bavovnacoin/node_controller/command_executor"
	"bavovnacoin/txo"
	"fmt"
)

func process() {
	blockchain.RestoreMempool()
	txo.RestoreCoinDatabase()
	InitAccountsData()
	blockchain.InitBlockchain()

	//go createTxRandom()
	for command_executor.Node_working {
		createAccoundRandom()
		addBlock()
	}
}

func Launch() {
	dbController.DB.OpenDb()
	isRpcStarted, err := networking.StartRPCListener()
	if !isRpcStarted {
		fmt.Println("Can't start RPC listener")
		fmt.Println(err)
	} else {
		fmt.Println("RPC listener started")
	}

	go process()
	node_controller.CommandHandler()
	blockchain.BackTransactionsToMempool()
	blockchain.WriteMempoolData()
	dbController.DB.CloseDb()
	account.WriteAccounts()
}
