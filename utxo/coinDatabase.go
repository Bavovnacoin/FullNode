package utxo

import (
	"bavovnacoin/byteArr"
	"fmt"
)

type UTXO struct {
	Address byteArr.ByteArr
	Sum     uint64
}

var UtxoList []UTXO

func DelFromUtxo(address byteArr.ByteArr, outind int) {
	ind := 0
	for i := 0; i < len(UtxoList); i++ {
		if UtxoList[i].Address.IsEqual(address) {
			if ind == outind {
				UtxoList = append(UtxoList[:i], UtxoList[i+1:]...)
				return
			}
			ind++
		}
	}
}

func AddToUtxo(address byteArr.ByteArr, sum uint64) {
	UtxoList = append(UtxoList, UTXO{Address: address, Sum: sum})
}

func ShowCoinDatabase() {
	for i := 0; i < len(UtxoList); i++ {
		fmt.Printf("Address: %s, sum: %d\n", UtxoList[i].Address.ToHexString(), UtxoList[i].Sum)
	}
}
