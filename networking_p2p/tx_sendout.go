package networking_p2p

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/byteArr"
	"bavovnacoin/hashing"
	"bavovnacoin/node/node_settings"
	"bavovnacoin/transaction"

	"github.com/libp2p/go-libp2p/core/peer"
)

func (pd *PeerData) TryHandleTx(data []byte, peerId peer.ID) bool {
	if data[0] == 8 {
		var tx transaction.Transaction
		isConv := byteArr.FromByteArr(data[1:], &tx)
		if !isConv {
			return false
		}

		var txHash byteArr.ByteArr
		txHash.SetFromHexString(hashing.SHA1(transaction.GetCatTxFields(tx)), 20)

		if !blockchain.IsTxInMempool(txHash) &&
			transaction.CheckTxMinFee(tx, node_settings.Settings.TxMinFee) &&
			blockchain.AddTxToMempool(tx, true) {
			pd.ProposeNewTx(tx)
			return true
		}
	}
	return false
}

func (pd *PeerData) ProposeNewTx(tx transaction.Transaction) bool {
	txByte, _ := byteArr.ToByteArr(tx)
	return pd.SendDataToAllConnectedPeers(append([]byte{8}, txByte...))
}
