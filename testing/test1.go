package testing

import (
	"bavovnacoin/account"
	"bavovnacoin/blockchain"
	"bavovnacoin/dbController"
	"bavovnacoin/node_controller"
	"bavovnacoin/node_controller/command_executor"
	"bavovnacoin/utxo"
)

func process() {
	blockchain.RestoreMempool()
	utxo.RestoreCoinDatabase()
	InitAccountsData()
	blockchain.InitBlockchain()

	go createTxRandom()
	for command_executor.Node_working {
		createAccoundRandom()
		addBlock()
	}
}

func Test1() {
	dbController.DB.OpenDb()
	go process()
	node_controller.CommandHandler()
	blockchain.BackTransactionsToMempool()
	blockchain.WriteMempoolData()
	dbController.DB.CloseDb()
	account.WriteAccounts()
}
