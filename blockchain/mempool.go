package blockchain

import (
	"bavovnacoin/byteArr"
	"bavovnacoin/hashing"
	"bavovnacoin/transaction"
	"log"
)

var Mempool []transaction.Transaction
var MempTxHashes = make(map[string]bool)
var MempInputHashes = make(map[string]bool)

func IsTxInMempool(txHash byteArr.ByteArr) bool {
	return MempTxHashes[txHash.ToHexString()]
}

func AreInputsInMempool(inputs []transaction.Input) bool {
	for _, inp := range inputs {
		if MempInputHashes[inp.GetHash()] {
			return false
		}
	}
	return true
}

func AddInputsToMempInpHashes(inputs []transaction.Input) {
	for _, inp := range inputs {
		MempInputHashes[inp.GetHash()] = true
	}
}

func RemInputsFromMempInpHashes(inputs []transaction.Input) {
	for _, inp := range inputs {
		delete(MempInputHashes, inp.GetHash())
	}
}

func AddTxToMempool(tx transaction.Transaction, allowVerify bool) bool {
	if allowVerify && !(transaction.VerifyTransaction(tx) || AreInputsInMempool(tx.Inputs)) {
		return false
	} else {
		insInd := findIndexSorted(transaction.GetTxFee(tx), tx.Locktime)

		if len(Mempool) != 0 {
			if insInd < len(Mempool) {
				Mempool = append(Mempool[:insInd+1], Mempool[insInd:]...)
				Mempool[insInd] = tx
				AddInputsToMempInpHashes(tx.Inputs)
				MempTxHashes[hashing.SHA1(transaction.GetCatTxFields(tx))] = true
				return true
			} else {
				Mempool = append(Mempool, tx)
				AddInputsToMempInpHashes(tx.Inputs)
				MempTxHashes[hashing.SHA1(transaction.GetCatTxFields(tx))] = true
				return true
			}
		} else {
			Mempool = append(Mempool, tx)
			AddInputsToMempInpHashes(tx.Inputs)
			MempTxHashes[hashing.SHA1(transaction.GetCatTxFields(tx))] = true
			return true
		}
	}
}

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

	for MempoolInd := 0; allSize < 1000000-coinbaseTxSize && MempoolInd < len(Mempool); MempoolInd++ {
		if Mempool[MempoolInd].Locktime < uint(BcLength) {
			allSize += transaction.ComputeTxSize(Mempool[MempoolInd])
			txForBlock = append(txForBlock, Mempool[MempoolInd])
		}
	}
	return txForBlock
}

func RemoveTxsFromMempool(txs []transaction.Transaction) {
	remCount := 0
	for mempInd := 0; mempInd < len(Mempool); mempInd++ {
		mempTxHash := hashing.SHA1(transaction.GetCatTxFields(Mempool[mempInd]))
		for _, txVal := range txs {
			currTxHash := hashing.SHA1(transaction.GetCatTxFields(txVal))
			if mempTxHash == currTxHash {
				Mempool = append(Mempool[:mempInd], Mempool[mempInd+1:]...)
				RemInputsFromMempInpHashes(txVal.Inputs)
				delete(MempTxHashes, hashing.SHA1(transaction.GetCatTxFields(txVal)))
				remCount++
				mempInd--
				break
			}
		}
		if remCount == len(txs) {
			break
		}
	}
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
			AddInputsToMempInpHashes(BlockForMining.Transactions[i].Inputs)
			MempTxHashes[hashing.SHA1(transaction.GetCatTxFields(BlockForMining.Transactions[i]))] = true
		}
		if len(BlockForMining.Transactions) > 1 {
			log.Printf("%d transactions are returned back to the mempool\n", len(BlockForMining.Transactions)-1)
		}
	}
}
