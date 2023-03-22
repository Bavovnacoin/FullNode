package networking

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/byteArr"
	"bavovnacoin/node_controller/node_settings"
	"bavovnacoin/transaction"
)

func (l *Listener) AddNewTxToMemp(txByteArr []byte, reply *Reply) error {
	var tx transaction.Transaction
	byteArr.FromByteArr(txByteArr, &tx)

	if transaction.CheckTxMinFee(tx, node_settings.Settings.TxMinFee) &&
		blockchain.AddTxToMempool(tx, true) {
		*reply = Reply{[]byte{1}}
	} else {
		*reply = Reply{[]byte{0}}
	}
	return nil
}
