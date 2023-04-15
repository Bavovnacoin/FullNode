package main

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/hashing"
	"fmt"
	"math/big"
)

func main() {
	//t := new(big.Int)
	//t.SetString("0000ffff00000000000000000000000000000000", 16)
	//fmt.Printf("%x", blockchain.TargetToBits(t))

	//println(big.NewInt(1).Exp(big.NewInt(2), big.NewInt(256), big.NewInt(0)).String()) // Chainwork length is not 32 bytes! Recalculate according to SHA1

	blockTarget := blockchain.BitsToTarget(0xffff12)
	fmt.Printf("%d", new(big.Int).Div(hashing.MaxNum, new(big.Int).Add(blockTarget, big.NewInt(1))))

	//node.Launch()

	//synchronization.GetCheckpHashes(2, 5)

}
