package main

import (
	"bavovnacoin/testing/singleFunctionTesting"
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
	} else if len(os.Args) == 2 && os.Args[1] == "-bft" { //TODO: fix it
		var bft singleFunctionTesting.BlockchainVerifTest
		bft.Launch(5, 3, 4)
	} else if len(os.Args) == 4 && os.Args[1] == "-cat" {
		testAmmount, err := strconv.Atoi(os.Args[3])

		if err != nil || testAmmount < 0 {
			println("Error: second parameter is anvalid")
			return
		}

		var cat singleFunctionTesting.CryptoTest
		if os.Args[2] == "ecdsa" || os.Args[2] == "sha1" {
			cat.Launch(os.Args[2], uint(testAmmount))
		}
	} else if len(os.Args) == 2 && os.Args[1] == "-dct" { //TODO: fix it
		var dct singleFunctionTesting.DifficultyChangeTest
		dct.Launch()
	} else if len(os.Args) == 3 && os.Args[1] == "-nct" {
		testAmmount, err := strconv.Atoi(os.Args[2])

		if err != nil || testAmmount < 0 {
			println("Error: second parameter is anvalid")
			return
		}

		var dct singleFunctionTesting.CommunicationTest
		dct.Launch(testAmmount)
	} else if len(os.Args) == 2 && os.Args[1] == "-pmt" { //TODO: fix it
		var pmt singleFunctionTesting.ParallelMiningTest
		pmt.Launch()
	}

	else {
		println("Unknown parameter")
	}

}
