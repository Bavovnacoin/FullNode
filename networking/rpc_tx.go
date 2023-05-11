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
		networking_p2p.Peer.ProposeNewTx(tx)
	} else {
		*reply = Reply{[]byte{0}}
	}
	return nil
}

func (c *Connection) SendTransaction(tx transaction.Transaction, isAccepted *bool) bool {
	byteArr, isConv := c.ToByteArr(tx)
	if !isConv {
		return false
	}

	var repl Reply
	err := c.client.Call("Listener.AddNewTxToMemp", byteArr, &repl)
	if err != nil {
		return false
	}

	if repl.Data[0] == 1 {
		*isAccepted = true
	} else {
		*isAccepted = false
	}
	return true
}
