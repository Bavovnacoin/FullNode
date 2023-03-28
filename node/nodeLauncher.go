package node

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/dbController"
	"bavovnacoin/networking"
	"bavovnacoin/node_controller"
	"bavovnacoin/node_controller/command_executor"
	"bavovnacoin/node_controller/node_settings"
	"bavovnacoin/synchronization"
	"bavovnacoin/txo"
	"fmt"
	"log"
	"runtime"
)

var NodeLaunched bool

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
	for command_executor.ComContr.FullNodeWorking {
		AddBlock(true)
	}
}

func LaunchFullNode() {
	command_executor.ComContr.FullNodeWorking = true
	dbController.DB.OpenDb()
	defer dbController.DB.CloseDb()
	StartRPC()
	defer networking.StopRPCListener()
	blockchain.InitBlockchain()
	txo.RestoreCoinDatabase()

	log.Println("Db synchronization...")
	syncRes := synchronization.StartSync(true, blockchain.BcLength)
	if !syncRes {
		input := ""
		for true {
			command_executor.ComContr.ClearConsole()
			log.Println("An error occured when synchronizing DB")
			log.Println("To continue enter \"Yes\". To back to the menu enter \"back\". ")
			fmt.Scan(&input)
			if input == "Yes" {
				break
			} else if input == "back" {
				return
			}
		}
	} else {
		log.Println("Db synchronization done")
	}

	if blockchain.BcLength == 0 {
		blockchain.FormGenesisBlock()
	}
	// panic(fmt.Sprintf("Bc len: %d", blockchain.BcLength))

	go NodeProcess()
	node_controller.CommandHandler()
}

func funcChoser(variant string) {
	if variant == "1" {
		command_executor.ComContr.ClearConsole()
		LaunchFullNode()
	} else if variant == "2" {
		node_settings.LaunchMenu(&node_settings.Settings)
	} else if variant == "3" {
		NodeLaunched = false
	}
}

func Launch() {
	NodeLaunched = true
	command_executor.ComContr.OpSys = runtime.GOOS
	node_settings.Settings.GetSettings()
	node_settings.Settings.InitSettingsValues()

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
