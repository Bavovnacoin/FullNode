package main

import (
	"bavovnacoin/node"
	"bavovnacoin/node/node_validator"
	"bavovnacoin/testing/loadTesting"
	"bavovnacoin/testing/loadTesting/simultNodeWork"
	"bavovnacoin/testing/singleFunctionTesting"
	"fmt"
	"os"
	"runtime"
	"strconv"
)

func main() {
	if len(os.Args) == 1 {
		println("Node launch")
		node.Launch()
	} else if len(os.Args) == 3 && os.Args[1] == "-snw" && os.Args[2] == "-l" {
		node_validator.LaunchValidatorNode()
	} else if len(os.Args) == 3 && os.Args[1] == "-snw" {
		nodesAmmount, err := strconv.Atoi(os.Args[2])

		if err != nil || nodesAmmount < 0 {
			println("Error: second parameter is anvalid")
			return
		}

		if nodesAmmount*5 > runtime.NumCPU() {
			fmt.Printf("Warning: ammount of nodes is higher than recomended (recomended: %d)\n", runtime.NumCPU()/5)
			var input string
			for true {
				println("Do you really want to continue [y/n]")
				fmt.Scanln(&input)

				if input == "y" {
					break
				} else if input == "n" {
					return
				}
			}
		}

		var snw simultNodeWork.SimultaneousNodeWork
		snw.Launch(nodesAmmount)
	} else if len(os.Args) == 5 && os.Args[1] == "-bvt" {
		blockAmmount, err := strconv.Atoi(os.Args[2])
		incblockAmmount, err := strconv.Atoi(os.Args[3])
		txAmmount, err := strconv.Atoi(os.Args[4])

		if err != nil || blockAmmount <= 0 || incblockAmmount <= 0 || txAmmount <= 0 {
			println("Error: one of the parameters is anvalid")
			return
		}

		var bft singleFunctionTesting.BlockchainVerifTest
		bft.Launch(blockAmmount, incblockAmmount, txAmmount)
	} else if len(os.Args) == 4 && os.Args[1] == "-cat" {
		testAmmount, err := strconv.Atoi(os.Args[3])

		if err != nil || testAmmount <= 0 {
			println("Error: second parameter is invalid")
			return
		}

		var cat singleFunctionTesting.CryptoTest
		if os.Args[2] == "ecdsa" || os.Args[2] == "sha1" {
			cat.Launch(os.Args[2], uint(testAmmount))
		}
	} else if len(os.Args) == 2 && os.Args[1] == "-dct" {
		var dct singleFunctionTesting.DifficultyChangeTest
		dct.Launch()
	} else if len(os.Args) == 3 && os.Args[1] == "-nct" {
		testAmmount, err := strconv.Atoi(os.Args[2])

		if err != nil || testAmmount <= 0 {
			println("Error: second parameter is invalid")
			return
		}

		var dct singleFunctionTesting.CommunicationTest
		dct.Launch(testAmmount)
	} else if len(os.Args) == 3 && os.Args[1] == "-pmt" {
		testAmmount, err := strconv.Atoi(os.Args[2])

		if err != nil || testAmmount <= 0 {
			println("Error: second parameter is invalid")
			return
		}

		var pmt singleFunctionTesting.ParallelMiningTest
		pmt.Launch(testAmmount)
	} else if len(os.Args) == 2 && os.Args[1] == "-rvt" {
		var pmt singleFunctionTesting.ReorganizationVerifTest
		pmt.Launch()
	} else if len(os.Args) == 5 && os.Args[1] == "-st" {
		mcBlocksAmmount, err := strconv.Atoi(os.Args[2])
		acBlocksAmmount, err := strconv.Atoi(os.Args[3])
		acAmmount, err := strconv.Atoi(os.Args[4])

		if err != nil || mcBlocksAmmount <= 0 || acBlocksAmmount <= 0 || acAmmount <= 0 {
			println("Error: one of the parameters is anvalid")
			return
		}
		var st singleFunctionTesting.SynchronizationTest
		st.Launch(mcBlocksAmmount, acBlocksAmmount, acAmmount)
	} else if len(os.Args) == 4 && os.Args[1] == "-tv" {
		corrTxAmmount, err := strconv.Atoi(os.Args[2])
		incTxAmmount, err := strconv.Atoi(os.Args[3])

		if err != nil || corrTxAmmount <= 0 || incTxAmmount <= 0 {
			println("Error: one of the parameters is anvalid")
			return
		}
		var tv singleFunctionTesting.TransVerifTest
		tv.TransactionsVerefication(corrTxAmmount, incTxAmmount)
	} else if len(os.Args) == 2 && os.Args[1] == "-ttv" {
		var ttv singleFunctionTesting.TransVerifTime
		ttv.TransVerifTime()
	} else if len(os.Args) == 3 && os.Args[1] == "-tmt" {
		mcAmmount, err := strconv.Atoi(os.Args[3])

		if err != nil || mcAmmount <= 0 {
			println("Error: one of the parameters is anvalid")
			return
		}
		var tmt singleFunctionTesting.TxoManagingTest
		tmt.Launch(uint64(mcAmmount))
	} else if len(os.Args) == 4 && os.Args[1] == "-mlt" {
		corrTxAmmount, err := strconv.Atoi(os.Args[2])
		incTxAmmount, err := strconv.Atoi(os.Args[3])

		if err != nil || corrTxAmmount <= 0 || incTxAmmount <= 0 {
			println("Error: one of the parameters is anvalid")
			return
		}
		var mlt loadTesting.MempoolLoadTest
		mlt.Launch(uint(corrTxAmmount), uint(incTxAmmount))
	}

	// else {
	// 	println("Unknown parameter")
	// }

}
