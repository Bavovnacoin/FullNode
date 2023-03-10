package txo

import (
	"bavovnacoin/byteArr"
	"bavovnacoin/dbController"
	"bavovnacoin/hashing"
	"fmt"
	"log"

	"github.com/syndtr/goleveldb/leveldb/util"
)

func SetUtxo(utxo TXO) bool {
	byteVal, isConv := byteArr.ToByteArr(utxo)
	if !isConv {
		return false
	}
	key := "utxo" + hashing.SHA1(utxo.OutTxHash.ToHexString()+fmt.Sprint(utxo.TxOutInd))
	return dbController.DB.SetValue(key, byteVal)
}

func GetUtxo(outTxHash byteArr.ByteArr, outInd int) (TXO, bool) {
	utxoByteArr, isValid := dbController.DB.GetValue("utxo" + hashing.SHA1(outTxHash.ToHexString()+fmt.Sprint(outInd)))
	var utxo TXO
	if isValid {
		byteArr.FromByteArr(utxoByteArr, &utxo)
		return utxo, true
	}
	return utxo, false
}

func RestoreCoinDatabase() bool {
	iter := dbController.DB.Db.NewIterator(util.BytesPrefix([]byte("utxo")), nil)
	for iter.Next() {
		var utxo TXO
		isConv := byteArr.FromByteArr(iter.Value(), &utxo)
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
	byteVal, isConv := byteArr.ToByteArr(txo)
	if !isConv {
		println("Problem with conversion txo")
		return false
	}
	key := "txo" + hashing.SHA1(txo.OutTxHash.ToHexString()+fmt.Sprint(txo.TxOutInd)+fmt.Sprint(txo.BlockHeight))
	return dbController.DB.SetValue(key, byteVal)
}

func IsOutAddrExist(addr byteArr.ByteArr) bool {
	iter := dbController.DB.Db.NewIterator(util.BytesPrefix([]byte("txo")), nil)
	for i := 0; iter.Next(); i++ {
		var txo TXO
		isConv := byteArr.FromByteArr(iter.Value(), &txo)
		if !isConv {
			return false
		}

		if txo.OutAddress.IsEqual(addr) {
			return true
		}
	}

	iter.Release()
	iter.Error()

	iter = dbController.DB.Db.NewIterator(util.BytesPrefix([]byte("utxo")), nil)
	for i := 0; iter.Next(); i++ {
		var txo TXO
		isConv := byteArr.FromByteArr(iter.Value(), &txo)
		if !isConv {
			return false
		}

		if txo.OutAddress.IsEqual(addr) {
			return true
		}
	}

	iter.Release()
	iter.Error()
	return false
}

func PrintSpentTxOuts() bool {
	i := 0
	iter := dbController.DB.Db.NewIterator(util.BytesPrefix([]byte("txo")), nil)
	for iter.Next() {
		var txo TXO
		isConv := byteArr.FromByteArr(iter.Value(), &txo)
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

func PrintUtxo() bool {
	i := 0
	iter := dbController.DB.Db.NewIterator(util.BytesPrefix([]byte("utxo")), nil)
	for iter.Next() {
		var txo TXO
		isConv := byteArr.FromByteArr(iter.Value(), &txo)
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
