package utxo

import (
	"bavovnacoin/byteArr"
	"fmt"
	"log"
)

type TXO struct {
	OutTxHash   byteArr.ByteArr
	TxOutInd    uint64
	Sum         uint64
	OutAddress  byteArr.ByteArr
	BlockHeight uint64
	IsSpent     bool
}

var CoinDatabase []TXO

func Spend(outTxHash byteArr.ByteArr, outind uint64) {
	for i := 0; i < len(CoinDatabase); i++ {
		if CoinDatabase[i].OutAddress.IsEqual(outTxHash) &&
			CoinDatabase[i].TxOutInd == outind {
			CoinDatabase[i].IsSpent = true
			return
		}
	}
}

func IsUtxoExists(txHash byteArr.ByteArr, outInd uint64) bool {
	for i := 0; i < len(CoinDatabase); i++ {
		if CoinDatabase[i].OutTxHash.IsEqual(txHash) && CoinDatabase[i].TxOutInd == outInd && !CoinDatabase[i].IsSpent {
			return true
		}
	}
	return false
}

func AddUtxo(outTxHash byteArr.ByteArr, txOutInd uint64, sum uint64,
	outAddress byteArr.ByteArr, blockHeight uint64) {
	CoinDatabase = append(CoinDatabase, TXO{OutTxHash: outTxHash, TxOutInd: txOutInd, Sum: sum, OutAddress: outAddress,
		BlockHeight: blockHeight, IsSpent: false})
}

func isSpentToStr(isSpent bool) string {
	if isSpent {
		return "Spent"
	}
	return "Unspent"
}

func PrintCoinDatabase() {
	log.Println("Utxo list:")
	for i := 0; i < len(CoinDatabase); i++ {
		fmt.Printf("[%d]. %s tx. Coins from transaction: %s (output num. %d) on address %s. Block height: %d. sum: %d\n",
			i, isSpentToStr(CoinDatabase[i].IsSpent), CoinDatabase[i].OutTxHash.ToHexString(),
			CoinDatabase[i].TxOutInd, CoinDatabase[i].OutAddress.ToHexString(), CoinDatabase[i].BlockHeight, CoinDatabase[i].Sum)
	}
}
