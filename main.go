package main

import "bavovnacoin/testing/singleFunctionTesting"

func main() {
	//t := new(big.Int)
	//t.SetString("0000ffff00000000000000000000000000000000", 16)
	//fmt.Printf("%x", blockchain.TargetToBits(t))

	var rv singleFunctionTesting.ReorganizationVerifTest
	rv.Launch()

	// var rv singleFunctionTesting.BlockchainVerifTest
	// rv.Launch(10, 2, 3)

	//node.Launch()

	//synchronization.GetCheckpHashes(2, 5)

}
