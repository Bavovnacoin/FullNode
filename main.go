package main

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/testing/singleFunctionTesting"
	"fmt"
	"math/big"
)

func main() {
	a, _ := new(big.Int).SetString("00ffff0000000000000000000000000000000000", 16)
	fmt.Printf("%x\n", blockchain.TargetToBits(a))
	var rv singleFunctionTesting.ReorganizationVerifTest
	rv.Launch()

	// node.Launch()
	//synchronization.GetCheckpHashes(1)

}
