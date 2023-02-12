package utxo

import (
	"bavovnacoin/byteArr"
	"bavovnacoin/dbController"
	"bavovnacoin/hashing"
	"fmt"
	"log"

	"github.com/syndtr/goleveldb/leveldb/util"
)

func SetUtxo(utxo TXO) bool {
	byteVal, isConv := dbController.ToByteArr(utxo)
	if !isConv {
		return false
	}
	key := "utxo" + hashing.SHA1(utxo.OutTxHash.ToHexString()+fmt.Sprint(utxo.TxOutInd)+fmt.Sprint(utxo.BlockHeight))
	return dbController.DB.SetValue(key, byteVal)
}

func RestoreCoinDatabase() bool {
	iter := dbController.DB.Db.NewIterator(util.BytesPrefix([]byte("utxo")), nil)
	for iter.Next() {
		var utxo TXO
		isConv := dbController.FromByteArr(iter.Value(), &utxo)
		if !isConv {
			log.Println("Problem when restoring coin database")
			return false
		}
		CoinDatabase = append(CoinDatabase, utxo)
	}

	iter.Release()
	err := iter.Error()
	if err != nil {
		log.Println("Problem when restoring coin database")
		return false
	}
	log.Println("Coin database restored")
	return true
}

func RemUtxo(OutTxHash byteArr.ByteArr, OutTxInd, blockHeight uint64) bool {
	key := "utxo" + hashing.SHA1(OutTxHash.ToHexString()+fmt.Sprint(OutTxInd)+fmt.Sprint(blockHeight))
	res := dbController.DB.RemoveValue(key)
	return res
}

func SetTxo(txo TXO) bool {
	byteVal, isConv := dbController.ToByteArr(txo)
	if !isConv {
		return false
	}
	key := "txo" + hashing.SHA1(txo.OutTxHash.ToHexString()+fmt.Sprint(txo.TxOutInd)+fmt.Sprint(txo.BlockHeight))
	return dbController.DB.SetValue(key, byteVal)
}

func PrintSpentTxOuts() bool {
	i := 0
	iter := dbController.DB.Db.NewIterator(util.BytesPrefix([]byte("txo")), nil)
	for iter.Next() {
		var txo TXO
		isConv := dbController.FromByteArr(iter.Value(), &txo)
		if !isConv {
			return false
		}
		txo.PrintTxo(i)
		i++
	}

	iter.Release()
	iter.Error()
	return true
}
