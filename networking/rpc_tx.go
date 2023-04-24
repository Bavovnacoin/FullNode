package networking

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/byteArr"
	"bavovnacoin/hashing"
	"bavovnacoin/node/node_controller/node_settings"
	"bavovnacoin/transaction"
)

func (l *Listener) AddNewTxToMemp(txByteArr []byte, reply *Reply) error {
	var tx transaction.Transaction
	byteArr.FromByteArr(txByteArr, &tx)

	if transaction.CheckTxMinFee(tx, node_settings.Settings.TxMinFee) &&
		blockchain.AddTxToMempool(tx, true) {
		*reply = Reply{[]byte{1}}
		ProposeTxToSettingsNodes(tx, "")
	} else {
		*reply = Reply{[]byte{0}}
	}
	return nil
}

func (l *Listener) AddProposedTxToMemp(txProposalByteArr []byte, reply *Reply) error {
	var txProp TxProposal
	byteArr.FromByteArr(txProposalByteArr, &txProp)

	if transaction.CheckTxMinFee(txProp.Tx, node_settings.Settings.TxMinFee) &&
		blockchain.AddTxToMempool(txProp.Tx, true) {
		*reply = Reply{[]byte{1}}
		ProposeTxToSettingsNodes(txProp.Tx, "")
	} else {
		*reply = Reply{[]byte{0}}
	}
	return nil
}

func (l *Listener) GetTxProposal(txByteArr []byte, reply *Reply) error {
	var txHash byteArr.ByteArr
	txHash.ByteArr = txByteArr
	if !blockchain.IsTxInMempool(txHash) {
		*reply = Reply{[]byte{1}}
	} else {
		*reply = Reply{[]byte{0}}
		return nil
	}
	return nil
}

func (c *Connection) ProposeTxToOtherNode(txHash []byte, tx transaction.Transaction) bool {
	var repl Reply
	err := c.client.Call("Listener.GetTxProposal", txHash, &repl)
	if err != nil {
		return false // Problem when accessing an RPC function
	}

	if repl.Data[0] == 1 {
		repl.Data = []byte{}

		var txProp TxProposal
		txProp.Tx = tx
		txProp.Address = node_settings.Settings.MyAddress
		propBytes, _ := c.ToByteArr(txProp)

		err := c.client.Call("Listener.AddProposedTxToMemp", propBytes, &repl)
		if err != nil || repl.Data[0] == 0 {
			return false // The node reverted this tx
		}
	} else {
		return false // The node is already has this tx
	}
	return true // No problems
}

func ProposeTxToSettingsNodes(tx transaction.Transaction, avoidAddress string) bool {
	var txHash byteArr.ByteArr
	txHashString := hashing.SHA1(transaction.GetCatTxFields(tx))
	txHash.SetFromHexString(txHashString, 20)

	var connection Connection
	var isNodesAccessible bool

	for i := 0; i < len(node_settings.Settings.OtherNodesAddresses); i++ {
		isNodesAccessible, i = connection.EstablishAddresses(node_settings.Settings.OtherNodesAddresses, i-1, avoidAddress)

		if !isNodesAccessible {
			return false
		}

		connection.ProposeTxToOtherNode(txHash.ByteArr, tx)
	}
	return true
}
