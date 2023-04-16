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
	"strings"
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
		genBlock := blockchain.FormGenesisBlock()
		networking.ProposeBlockToSettingsNodes(genBlock, "")
	}

	go NodeProcess()
	node_controller.CommandHandler()
}

func funcChoser(variant string, isNodeLaunchAllowed bool) {
	if variant == "1" && isNodeLaunchAllowed {
		command_executor.ComContr.ClearConsole()
		LaunchFullNode()
	} else if variant == "2" && isNodeLaunchAllowed || variant == "1" && !isNodeLaunchAllowed {
		node_settings.LaunchMenu(&node_settings.Settings)
	} else if variant == "3" && isNodeLaunchAllowed || variant == "2" && !isNodeLaunchAllowed {
		NodeLaunched = false
	}
}

func getNodeLaunchSettingsError() string {
	var errMess []string
	if node_settings.Settings.RewardAddress == "" {
		errMess = append(errMess, "reward address (1-6)")
	}
	if node_settings.Settings.MyAddress == "" {
		errMess = append(errMess, "node address (1-5)")
	}
	return strings.Join(errMess, ", ")
}

func Launch() {
	NodeLaunched = true
	command_executor.ComContr.OpSys = runtime.GOOS
	node_settings.Settings.GetSettings()
	node_settings.Settings.InitSettingsValues()

	//os.RemoveAll("data") // TODO: remove!

	var variant string
	for NodeLaunched {
		command_executor.ComContr.ClearConsole()
		println("Choose a variant and press the right button")
		launchMesErr := getNodeLaunchSettingsError()
		if launchMesErr != "" {
			fmt.Printf("Can't start a node. You need to manage: %s\n", launchMesErr)
		}

		var btn int = 1
		if launchMesErr == "" {
			fmt.Printf("%d. Launch node\n", btn)
			btn++
		}
		fmt.Printf("%d. Manage settings\n", btn)
		fmt.Printf("%d. Exit\n", btn+1)

		fmt.Scan(&variant)
		funcChoser(variant, launchMesErr == "")
	}

	command_executor.ComContr.ClearConsole()
	println("Thank you for supporting Bavovnacoin network!")
}
