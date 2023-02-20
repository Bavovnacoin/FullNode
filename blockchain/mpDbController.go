package blockchain

import (
	"bavovnacoin/dbController"
	"bavovnacoin/transaction"
	"fmt"
	"log"
)

func SetMempLen(len uint64) bool {
	byteVal, isConv := dbController.ToByteArr(len)
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
	isConv := dbController.FromByteArr(value, &len)
	if !isConv {
		return 0, false
	}
	return len, true
}

func WriteTxToMempool(id uint64, tx transaction.Transaction) bool {
	byteVal, isConv := dbController.ToByteArr(tx)
	if !isConv {
		return false
	}

	return dbController.DB.SetValue("mp"+fmt.Sprint(id), byteVal)
}

func GetTxFromMempool(id uint64) (transaction.Transaction, bool) {
	var tx transaction.Transaction
	value, isGotten := dbController.DB.GetValue("mp" + fmt.Sprint(id))
	if !isGotten {
		return tx, false
	}

	isConv := dbController.FromByteArr(value, &tx)
	if !isConv {
		return tx, false
	}

	return tx, true
}

func WriteMempoolData() {
	oldLen, _ := GetMempLen()
	var i uint64
	for ; i < uint64(len(Mempool)); i++ {
		WriteTxToMempool(i, Mempool[i])
	}

	if oldLen > uint64(len(Mempool)) {
		RemoveMPRange(i, oldLen)
	}
	SetMempLen(uint64(len(Mempool)))
}

func RemoveMPRange(start, end uint64) {
	for ; start < end; start++ {
		res := dbController.DB.RemoveValue("mp" + fmt.Sprint(start))
		if res {
			break
		}
	}
}

func RestoreMempool() {
	mempLen, _ := GetMempLen()
	var i uint64
	for ; i < mempLen; i++ {
		tx, val := GetTxFromMempool(i)
		if !val {
			log.Println("Problem when restoring mp")
			break
		}
		Mempool = append(Mempool, tx)
	}
	if mempLen != 0 {
		log.Println("Mempool restored")
	} else {
		log.Println("Mempool is empty")
	}
}
