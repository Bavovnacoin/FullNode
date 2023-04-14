package main

import loadtesting "bavovnacoin/testing/loadTesting"

func main() {
	//t := new(big.Int)
	//t.SetString("0000ffff00000000000000000000000000000000", 16)
	//fmt.Printf("%x", blockchain.TargetToBits(t))

	//node.Launch()

	var lt loadtesting.LoadTest
	lt.StartLoadTest(10, 10)

	//synchronization.GetCheckpHashes(2, 5)

}
