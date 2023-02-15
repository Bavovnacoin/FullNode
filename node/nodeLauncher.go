package node

import (
	"bavovnacoin/account"
	"bavovnacoin/blockchain"
	"bavovnacoin/dbController"
	"bavovnacoin/node_controller"
	"bavovnacoin/node_controller/command_executor"
	"bavovnacoin/txo"
)

func process() {
	blockchain.RestoreMempool()
	txo.RestoreCoinDatabase()
	InitAccountsData()
	blockchain.InitBlockchain()

	go createTxRandom()
	for command_executor.Node_working {
		createAccoundRandom()
		addBlock()
	}
}

func Launch() {
	dbController.DB.OpenDb()
	go process()
	node_controller.CommandHandler()
	blockchain.BackTransactionsToMempool()
	blockchain.WriteMempoolData()
	dbController.DB.CloseDb()
	account.WriteAccounts()
}
