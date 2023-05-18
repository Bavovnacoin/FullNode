package main

import (
	"bavovnacoin/node"
	"bavovnacoin/node/node_validator"
	"bavovnacoin/testing/loadTesting/simultNodeWork"
	"os"
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

		// if nodesAmmount*5 > runtime.NumCPU() {
		// 	fmt.Printf("Warning: ammount of nodes is higher than recomended (recomended: %d)\n", runtime.NumCPU()/5)
		// 	var input string
		// 	for true {
		// 		println("Do you really want to continue [y/n]")
		// 		fmt.Scanln(&input)

		// 		if input == "y" {
		// 			break
		// 		} else if input == "n" {
		// 			return
		// 		}
		// 	}
		// }

		var snw simultNodeWork.SimultaneousNodeWork
		snw.Launch(nodesAmmount)
	} else {
		println("Unknown parameter")
	}

}
