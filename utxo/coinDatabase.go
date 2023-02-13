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
}

var CoinDatabase []TXO

func Spend(outTxHash byteArr.ByteArr, outind uint64) {
	var txoToSpend TXO

	for i := 0; i < len(CoinDatabase); i++ {
		if CoinDatabase[i].OutTxHash.IsEqual(outTxHash) &&
			CoinDatabase[i].TxOutInd == outind {
			txoToSpend = CoinDatabase[i]
			CoinDatabase = append(CoinDatabase[:i], CoinDatabase[i+1:]...)
			log.Println("Utxo removed array")
			remRes := RemUtxo(outTxHash, outind, txoToSpend.TxOutInd)
			if remRes {
				log.Println("Utxo removed db")
			}
			break
		}
	}
	setRes := SetTxo(txoToSpend)
	if setRes {
		log.Println("Utxo has been added to db")
	}
}

func IsUtxoExists(txHash byteArr.ByteArr, outInd uint64) bool {
	for i := 0; i < len(CoinDatabase); i++ {
		if CoinDatabase[i].OutTxHash.IsEqual(txHash) && CoinDatabase[i].TxOutInd == outInd {
			return true
		}
	}
	return false
}

func AddUtxo(outTxHash byteArr.ByteArr, txOutInd uint64, sum uint64,
	outAddress byteArr.ByteArr, blockHeight uint64) {
	utxo := TXO{OutTxHash: outTxHash, TxOutInd: txOutInd, Sum: sum, OutAddress: outAddress,
		BlockHeight: blockHeight}
	res := SetUtxo(utxo)
	if !res {
		log.Println("Problem when adding utxo to database")
	}
	CoinDatabase = append(CoinDatabase, utxo)
}

func (txo *TXO) PrintTxo(i int) {
	fmt.Printf("[%d]. Coins from transaction: %s (output num. %d) on address %s. Block height: %d. sum: %d\n",
		i, txo.OutTxHash.ToHexString(), txo.TxOutInd, txo.OutAddress.ToHexString(), txo.BlockHeight, txo.Sum)
}

// TODO: print unspent and spent outputs
func PrintCoinDatabase() {
	log.Println("Utxo list:")
	for i := 0; i < len(CoinDatabase); i++ {
		CoinDatabase[i].PrintTxo(i)
	}
	log.Println("Txo list:")
	PrintSpentTxOuts()
}
