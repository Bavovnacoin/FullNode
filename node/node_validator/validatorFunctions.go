package node_validator

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/dbController"
	"bavovnacoin/networking"
	"bavovnacoin/networking_p2p"
	"bavovnacoin/node/node_controller"
	"bavovnacoin/node/node_controller/command_executor"
	"bavovnacoin/synchronization"
	"bavovnacoin/txo"
	"fmt"
	"log"
)

func NodeProcess() {
	for command_executor.ComContr.FullNodeWorking {
		AddBlock(true)
	}
}

func BlockGen(allowCommandHandler bool) {
	if blockchain.BcLength == 0 {
		genBlock := blockchain.FormGenesisBlock()
		networking.ProposeBlockToSettingsNodes(genBlock, "")
	}

	if allowCommandHandler {
		go NodeProcess()
		node_controller.CommandHandler()
	} else {
		NodeProcess()
	}
}

func StartRPC() {
	isRpcStarted, err := networking.StartRPCListener()
	if !isRpcStarted {
		fmt.Println("Can't start RPC listener")
		fmt.Println(err)
	} else {
		fmt.Println("RPC listener started")
	}
}

func LaunchValidatorNode() {
	command_executor.ComContr.FullNodeWorking = true
	dbController.DB.OpenDb()
	defer dbController.DB.CloseDb()
	networking_p2p.StartP2PCommunication()
	//StartRPC()
	defer networking.StopRPCListener()
	blockchain.InitBlockchain()
	txo.RestoreCoinDatabase()

	log.Println("Db synchronization...")
	syncRes := synchronization.StartSync(true, blockchain.BcLength)
	if !syncRes {
		input := ""
		for true {
			//command_executor.ComContr.ClearConsole()
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

	BlockGen(true)
}
