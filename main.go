package main

import (
	"bavovnacoin/node"
	"fmt"
)

func ByteToString(arr []byte) string {
	str := ""
	for i := 0; i < len(arr); i++ {
		str += fmt.Sprint(arr[i])
	}
	return str
}

func main() {
	// a, _ := new(big.Int).SetString("00ffff0000000000000000000000000000000000", 16)
	// fmt.Printf("%x\n", blockchain.TargetToBits(a))

	node.Launch()
}
