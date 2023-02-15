package blockchain

import (
	"bavovnacoin/byteArr"
	"bavovnacoin/transaction"
	"log"
)

var Mempool []transaction.Transaction

func ValidateTransaction(tx transaction.Transaction) bool {
	if !transaction.VerifyTransaction(tx) {
		return false
	}

	// TODO: Remove!
	for j := 0; j < len(tx.Inputs); j++ {
		for i := 0; i < len(Mempool); i++ { // Check same input in mempool (TODO: find more effective way)
			for k := 0; k < len(Mempool[i].Inputs); k++ {
				if Mempool[i].Inputs[k].TxHash.IsEqual(tx.Inputs[j].TxHash) &&
					Mempool[i].Inputs[k].OutInd == tx.Inputs[j].OutInd {
					return false
				}
			}
		}
	}
	return true
}

func AddTxToMempool(tx transaction.Transaction, allowValidate bool) bool {
	if allowValidate && !ValidateTransaction(tx) {
		return false
	} else {
		fee := transaction.GetTxFee(tx)
		insInd := findIndexSorted(fee, tx.Locktime)

		if len(Mempool) != 0 {
			if insInd < len(Mempool) {
				Mempool = append(Mempool[:insInd+1], Mempool[insInd:]...)
				Mempool[insInd] = tx
				return true
			} else {
				Mempool = append(Mempool, tx)
				return true
			}
		} else {
			Mempool = append(Mempool, tx)
			return true
		}
	}
}

// Make binary search???
func findIndexSorted(fee uint64, locktime uint) int {
	for i := 0; i < len(Mempool); i++ {
		txFee := transaction.GetTxFee(Mempool[i])
		if txFee == fee {
			if Mempool[i].Locktime < locktime {
				return i
			}
		}
		if txFee < fee {
			return i
		}
	}
	return len(Mempool)
}

func GetTransactionsFromMempool(coinbaseTxSize int) []transaction.Transaction {
	var txForBlock []transaction.Transaction
	allSize := 0
	MempoolInd := 0

	for allSize < 1000000-coinbaseTxSize && MempoolInd < len(Mempool) {
		allSize += transaction.ComputeTxSize(Mempool[MempoolInd])

		if !transaction.VerifyTransaction(Mempool[MempoolInd]) {
			Mempool = append(Mempool[:MempoolInd], Mempool[MempoolInd+1:]...)
			log.Println("Deleted wrong transaction from mempool.")
		} else if Mempool[MempoolInd].Locktime < uint(BcLength) {
			txForBlock = append(txForBlock, Mempool[MempoolInd])
			Mempool = append(Mempool[:MempoolInd], Mempool[MempoolInd+1:]...)
		} else {
			MempoolInd++
		}
	}
	return txForBlock
}

func IsAddressInMempool(address byteArr.ByteArr) bool {
	for i := 0; i < len(Mempool); i++ {
		for j := 0; j < len(Mempool[i].Inputs); j++ {
			if Mempool[i].Inputs[j].TxHash.IsEqual(address) {
				return true
			}
		}
	}
	return false
}

func PrintMempool() {
	mempoolMes := "Mempool:"
	if len(Mempool) == 0 {
		mempoolMes += " empty"
	}
	log.Println(mempoolMes)
	for i := 0; i < len(Mempool); i++ {
		log.Printf("[%d]. Fee: %d, coins: %d, locktime: %d\n", i, transaction.GetTxFee(Mempool[i]),
			transaction.GetOutputValue(Mempool[i].Outputs), Mempool[i].Locktime)
	}
}

func BackTransactionsToMempool() {
	if IsMempAdded {
		for i := 1; i < len(BlockForMining.Transactions); i++ {
			AddTxToMempool(BlockForMining.Transactions[i], false)
		}
		if len(BlockForMining.Transactions) > 1 {
			log.Printf("%d transactions are returned back to mempool\n", len(BlockForMining.Transactions)-1)
		}
	}
}
