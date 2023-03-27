package blockchain

import (
	"bavovnacoin/byteArr"
	"bavovnacoin/dbController"
	"bavovnacoin/hashing"
	"bavovnacoin/transaction"
	"fmt"
)

func SetMempLen(len uint64) bool {
	byteVal, isConv := byteArr.ToByteArr(len)
	if !isConv {
		return false
	}
	return dbController.DB.SetValue("mpLength", byteVal)
}

func GetMempLen() (uint64, bool) {
	value, isGotten := dbController.DB.GetValue("mpLength")
	if !isGotten {
		return 0, false
	}

	var len uint64
	isConv := byteArr.FromByteArr(value, &len)
	if !isConv {
		return 0, false
	}
	return len, true
}

func WriteTxToMempool(tx transaction.Transaction) bool {
	byteVal, isConv := byteArr.ToByteArr(tx)
	if !isConv {
		return false
	}

	return dbController.DB.SetValue("mp"+hashing.SHA1(transaction.GetCatTxFields(tx)), byteVal)
}

// TODO: get tx by hash
func GetTxFromMempool(txHash byteArr.ByteArr) (transaction.Transaction, bool) {
	var tx transaction.Transaction
	value, isGotten := dbController.DB.GetValue("mp" + txHash.ToHexString())
	if !isGotten {
		return tx, false
	}

	isConv := byteArr.FromByteArr(value, &tx)
	if !isConv {
		return tx, false
	}

	return tx, true
}

func RemoveMPRange(start, end uint64) {
	for ; start < end; start++ {
		res := dbController.DB.RemoveValue("mp" + fmt.Sprint(start))
		if res {
			break
		}
	}
}

// func WriteMempoolData() {
// 	oldLen, _ := GetMempLen()
// 	var i uint64
// 	for ; i < uint64(len(Mempool)); i++ {
// 		WriteTxToMempool(i, Mempool[i])
// 	}

// 	if oldLen > uint64(len(Mempool)) {
// 		RemoveMPRange(i, oldLen)
// 	}
// 	SetMempLen(uint64(len(Mempool)))
// }

// func RestoreMempool() {
// 	mempLen, _ := GetMempLen()
// 	var i uint64
// 	for ; i < mempLen; i++ {
// 		tx, val := GetTxFromMempool(i)
// 		if !val {
// 			log.Println("Problem when restoring mp")
// 			break
// 		}
// 		Mempool = append(Mempool, tx)
// 	}
// 	if mempLen != 0 {
// 		log.Println("Mempool restored")
// 	} else {
// 		log.Println("Mempool is empty")
// 	}
// }
