package networking

import (
	"bavovnacoin/byteArr"
	"bavovnacoin/txo"
)

func (l *Listener) GetUtxoByAddr(addressesByte []byte, reply *Reply) {
	var addresses []byteArr.ByteArr
	println(byteArr.FromByteArr(addressesByte, &addresses))

	var addrUtxo []txo.TXO
	for i := len(txo.CoinDatabase); i > 0; i-- {
		for j := len(addresses); j > 0; j-- {
			if txo.CoinDatabase[i].OutAddress.IsEqual(addresses[j]) {
				addrUtxo = append(addrUtxo, txo.CoinDatabase[i])
			}
		}
	}

	byteAddrUtxo, res := byteArr.ToByteArr(addrUtxo)
	if !res {
		byteAddrUtxo = []byte("false")
	}
	*reply = Reply{byteAddrUtxo}
}
