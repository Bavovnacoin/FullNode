package node

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/dbController"
	"bavovnacoin/networking"
	"bavovnacoin/node_controller"
	"bavovnacoin/node_controller/command_executor"
	"bavovnacoin/node_controller/node_settings"
	"bavovnacoin/txo"
	"fmt"
	"runtime"
)

var NodeLaunched bool
var NodeSettings node_settings.NodeSettings

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

	for command_executor.ComContr.FullNodeWorking {
		AddBlock(true)
	}
}

func LaunchFullNode() {
	command_executor.ComContr.FullNodeWorking = true
	dbController.DB.OpenDb()
	StartRPC()
	go NodeProcess()
	node_controller.CommandHandler()
	blockchain.BackTransactionsToMempool()
	blockchain.WriteMempoolData()
	dbController.DB.CloseDb()
}

func funcChoser(variant string) {
	if variant == "1" {
		command_executor.ComContr.ClearConsole()
		LaunchFullNode()
	} else if variant == "2" {
		node_settings.LaunchMenu(&NodeSettings)
	} else if variant == "3" {
		NodeLaunched = false
	}
}

func Launch() {
	NodeLaunched = true
	command_executor.ComContr.OpSys = runtime.GOOS
	NodeSettings.GetSettings()
	NodeSettings.InitSettingsValues()

	var variant string
	for NodeLaunched {
		command_executor.ComContr.ClearConsole()
		println("Choose a variant and press the right button")
		println("1. Launch node")
		println("2. Manage settings")
		println("3. Exit")

		fmt.Scan(&variant)
		funcChoser(variant)
	}

	command_executor.ComContr.ClearConsole()
	println("Thank you for supporting Bavovnacoin network!")
}
