package node_audithor

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/dbController"
	"bavovnacoin/networking"
	"bavovnacoin/networking_p2p"
	"bavovnacoin/node/node_controller"
	"bavovnacoin/node/node_controller/command_executor"
	"bavovnacoin/txo"
	"fmt"
	"log"
	"time"
)

func StartRPC() {
	var isRpcStarted bool
	var err error
	networking.Inbound, isRpcStarted, err = networking.StartRPCListener()

	if !isRpcStarted {
		fmt.Println("Can't start RPC listener")
		fmt.Println(err)
	} else {
		fmt.Println("RPC listener started")
	}
}

func NodeProcess() {
	for command_executor.ComContr.FullNodeWorking {
		blocks, res := RecieveBlocks()
		if res == 1 {
			println("Can not connect to any address. Please, manage your addresses in the settings")
			time.Sleep(time.Second * 10)
		} else {
			ReorgTests(blocks)
		}

	}
}

func MakeAudit() {
	GetReorgData()
	go NodeProcess()
	node_controller.CommandHandler()
}

func LaunchAudithor() {
	command_executor.ComContr.FullNodeWorking = true
	dbController.DbPath = "data/AudNode"
	dbController.DB.OpenDb()
	defer dbController.DB.CloseDb()
	//StartRPC()
	defer networking.StopRPCListener()
	blockchain.InitBlockchain()
	txo.RestoreCoinDatabase()

	log.Println("Db synchronization...")
	syncRes := networking_p2p.Peer.StartSync()
	for !networking_p2p.IsSyncEnded {
		time.Sleep(20 * time.Millisecond)
	}
	if !syncRes {
		input := ""
		for true {
			command_executor.ComContr.ClearConsole()
			log.Println("An error occured when synchronizing DB")
			log.Println("To back to the menu enter \"back\". ")
			fmt.Scan(&input)
			if input == "back" {
				return
			}
		}
	} else {
		log.Println("Db synchronization done")
	}

	MakeAudit()
	//BlockGen(true)
}
