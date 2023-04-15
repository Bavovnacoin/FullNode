package main

import "bavovnacoin/node"

func main() {
	//t := new(big.Int)
	//t.SetString("0000ffff00000000000000000000000000000000", 16)
	//fmt.Printf("%x", blockchain.TargetToBits(t))

	//println(big.NewInt(1).Exp(big.NewInt(2), big.NewInt(256), big.NewInt(0)).String()) // Chainwork length is not 32 bytes! Recalculate according to SHA1

	node.Launch()

	//synchronization.GetCheckpHashes(2, 5)

}
