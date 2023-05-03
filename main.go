package main

import (
	"bavovnacoin/ecdsa"
	"bavovnacoin/hashing"
)

func main() {
	// a, _ := new(big.Int).SetString("00ffff0000000000000000000000000000000000", 16)
	// fmt.Printf("%x\n", blockchain.TargetToBits(a))
	// var rv singleFunctionTesting.ReorganizationVerifTest
	// rv.Launch()

	ecdsa.InitValues()
	kp := ecdsa.GenKeyPair()
	mes := hashing.SHA1("abcd")

	s := ecdsa.Sign(mes, kp.PrivKey)
	println(ecdsa.Verify(kp.PublKey, s, mes))

	// node.Launch()
	//synchronization.GetCheckpHashes(1)

}
