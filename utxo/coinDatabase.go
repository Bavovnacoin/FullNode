package utxo

import (
	"bavovnacoin/address"
	"fmt"
)

type UTXO struct {
	Address address.Address
	Sum     uint64
}

var UtxoList []UTXO

func DelFromUtxo(address address.Address, outind int) {
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

func AddToUtxo(address address.Address, sum uint64) {
	UtxoList = append(UtxoList, UTXO{Address: address, Sum: sum})
}

func ShowCoinDatabase() {
	for i := 0; i < len(UtxoList); i++ {
		fmt.Printf("Address: %s, sum: %d\n", UtxoList[i].Address.ToHexString(), UtxoList[i].Sum)
	}
}
