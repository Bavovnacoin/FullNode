package main

import "bavovnacoin/node"

func main() {
	//t := new(big.Int)
	//t.SetString("0000ffff00000000000000000000000000000000", 16)
	//fmt.Printf("%x", blockchain.TargetToBits(t))

	node.Launch()

	//synchronization.GetCheckpHashes(2, 5)

	// d6be9343c40bbedaa647ac81763bedf8ca3177f9
	//[100 54 98 101 57 51 52 51 99 52 48 98 98 101 100 97 97 54 52 55 97 99 56 49 55 54 51 98 101 100 102 56 99 97 51 49 55 55 102 57]
	//h, _ := new(big.Int).SetString("d6be9343c40bbedaa647ac81763bedf8ca3177f9", 16)
	//h := new(big.Int).SetBytes([]byte{100, 54, 98, 101, 57, 51, 52, 51, 99, 52, 48, 98, 98, 101, 100, 97, 97, 54, 52, 55, 97, 99, 56, 49, 55, 54, 51, 98, 101, 100, 102, 56, 99, 97, 51, 49, 55, 55, 102, 57})
	//fmt.Println(h.Bytes())
	//println(string(byte(100)))
}
