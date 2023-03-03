package networking

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/byteArr"
	"bavovnacoin/transaction"
)

func (l *Listener) AddNewTxToMemp(txByteArr []byte, reply *Reply) error {
	var tx transaction.Transaction
	byteArr.FromByteArr(txByteArr, &tx)
	isAdded := blockchain.AddTxToMempool(tx, true)
	if isAdded {
		*reply = Reply{[]byte{1}}
	} else {
		*reply = Reply{[]byte{0}}
	}
	return nil
}
