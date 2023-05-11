package main

import (
	loadtesting "bavovnacoin/testing/loadTesting"
)

func main() {
	// a, _ := new(big.Int).SetString("00ffff0000000000000000000000000000000000", 16)
	// fmt.Printf("%x\n", blockchain.TargetToBits(a))

	var mlt loadtesting.MempoolLoadTest
	mlt.Launch()

	// node.Launch()
	//synchronization.GetCheckpHashes(1)

}
