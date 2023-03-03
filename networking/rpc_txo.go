package networking

import (
	"bavovnacoin/byteArr"
	"bavovnacoin/txo"
)

func (l *Listener) GetUtxoByAddr(addressesByte []byte, reply *Reply) error {
	var addresses []byteArr.ByteArr
	byteArr.FromByteArr(addressesByte, &addresses)

	var addrUtxo []txo.TXO
	for i := len(txo.CoinDatabase) - 1; i >= 0; i-- {
		for j := len(addresses) - 1; j >= 0; j-- {
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
	return nil
}
