package networking

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/byteArr"
	"bavovnacoin/networking_p2p"
	"bavovnacoin/node/node_controller/node_settings"
	"bavovnacoin/transaction"
)

func (l *Listener) AddNewTxToMemp(txByteArr []byte, reply *Reply) error {
	var tx transaction.Transaction
	byteArr.FromByteArr(txByteArr, &tx)

	if transaction.CheckTxMinFee(tx, node_settings.Settings.TxMinFee) &&
		blockchain.AddTxToMempool(tx, true) {
		*reply = Reply{[]byte{1}}
		networking_p2p.ProposeNewTx(tx)
	} else {
		*reply = Reply{[]byte{0}}
	}
	return nil
}
