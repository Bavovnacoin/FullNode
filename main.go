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

	// ecdsa.InitValues()
	// data, _ := new(big.Int).SetString(ecdsa.GenPrivKey(), 16)
	// privKey := data.Bytes()
	// pkString := hex.EncodeToString(privKey)
	// fmt.Println(privKey)
	// cr := cryption.AES_encrypt(pkString, "a")
	// fmt.Println(cr)
	// dcr := cryption.AES_decrypt(cr, "a")
	// fmt.Println(hex.DecodeString(dcr))

	// fmt.Println(fmt.Sprintf("%s", privKey))
	// cr := cryption.AES_encrypt(string(privKey), "a")
	// fmt.Println(cryption.AES_decrypt(cr, "a"))
}
