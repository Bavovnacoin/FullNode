package main

import (
	"bavovnacoin/blockchain"
	"fmt"
	"math/big"
)

func main() {
	a, _ := new(big.Int).SetString("000ffff000000000000000000000000000000000", 16)
	println(fmt.Sprintf("%x", blockchain.TargetToBits(a)))
}
