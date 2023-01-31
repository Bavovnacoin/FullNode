package blockchain

import (
	"bavovnacoin/transaction"
	"bavovnacoin/utxo"
	"log"
)

var Mempool []transaction.Transaction

func ValidateTransaction(tx transaction.Transaction) bool {
	if !transaction.VerifyTransaction(tx) || tx.Locktime < uint(len(Blockchain)) && tx.Locktime != 0 {
		return false
	}

	for j := 0; j < len(tx.Inputs); j++ {
		for i := 0; i < len(Mempool); i++ { // Check same input in mempool (TODO: find more effective way)
			for k := 0; k < len(Mempool[i].Inputs); k++ {
				if Mempool[i].Inputs[k].HashAdr == tx.Inputs[j].HashAdr &&
					Mempool[i].Inputs[k].OutInd == tx.Inputs[j].OutInd { // hash address and ind
					return false
				}
			}
		}

		isExist := false
		for i := 0; i < len(utxo.UtxoList); i++ {
			if utxo.UtxoList[i].Address == tx.Inputs[j].HashAdr {
				isExist = true
			}
		}

		if !isExist {
			return false
		}
	}
	return true
}

func AddTxToMempool(tx transaction.Transaction) bool {
	if ValidateTransaction(tx) {
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
	return false
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
		} else if Mempool[MempoolInd].Locktime < uint(len(Mempool)) {
			txForBlock = append(txForBlock, Mempool[MempoolInd])
			Mempool = append(Mempool[:MempoolInd], Mempool[MempoolInd+1:]...)
		} else {
			MempoolInd++
		}
	}
	return txForBlock
}

func IsAddressInMempool(address string) bool {
	for i := 0; i < len(Mempool); i++ {
		for j := 0; j < len(Mempool[i].Inputs); j++ {
			if Mempool[i].Inputs[j].HashAdr == address {
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
			transaction.GetOutputSum(Mempool[i].Outputs), Mempool[i].Locktime)
	}
}
