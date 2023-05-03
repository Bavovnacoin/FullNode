package txo

import (
	"bavovnacoin/byteArr"
	"bavovnacoin/dbController"
	"fmt"
	"log"

	"github.com/syndtr/goleveldb/leveldb/util"
)

func SetUtxo(utxo TXO) bool {
	byteVal, isConv := byteArr.ToByteArr(utxo)
	if !isConv {
		return false
	}
	key := "utxo" + utxo.OutTxHash.ToHexString() + fmt.Sprint(utxo.TxOutInd) + fmt.Sprint(utxo.BlockHeight)
	return dbController.DB.SetValue(key, byteVal)
}

func GetUtxos(outTxHash byteArr.ByteArr, outInd int) ([]TXO, bool) {
	var utxoRes []TXO
	iter := dbController.DB.Db.NewIterator(util.BytesPrefix([]byte("utxo"+outTxHash.ToHexString()+fmt.Sprint(outInd))), nil)
	for iter.Next() {
		var utxo TXO
		isConv := byteArr.FromByteArr(iter.Value(), &utxo)
		if !isConv {
			return utxoRes, false
		}
		utxoRes = append(utxoRes, utxo)
	}

	iter.Release()
	err := iter.Error()
	if err != nil {
		return utxoRes, false
	}

	return utxoRes, true
}

func RemUtxo(OutTxHash byteArr.ByteArr, OutTxInd, blockHeight uint64) bool {
	key := "utxo" + OutTxHash.ToHexString() + fmt.Sprint(OutTxInd) + fmt.Sprint(blockHeight)
	res := dbController.DB.RemoveValue(key)
	return res
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

func SetTxo(txo TXO) bool {
	byteVal, isConv := byteArr.ToByteArr(txo)
	if !isConv {
		return false
	}
	key := "txo" + txo.OutTxHash.ToHexString() + fmt.Sprint(txo.TxOutInd) + fmt.Sprint(txo.BlockHeight)
	return dbController.DB.SetValue(key, byteVal)
}

func GetTxo(outTxHash byteArr.ByteArr, outInd int, BlockHeight uint64) (TXO, bool) {
	utxoByteArr, isValid := dbController.DB.GetValue("txo" + outTxHash.ToHexString() + fmt.Sprint(outInd) + fmt.Sprint(BlockHeight))
	var txo TXO
	if isValid {
		byteArr.FromByteArr(utxoByteArr, &txo)
		return txo, true
	}
	return txo, false
}

func GetTxos(outTxHash byteArr.ByteArr, outInd int) ([]TXO, bool) {
	var txoRes []TXO
	iter := dbController.DB.Db.NewIterator(util.BytesPrefix([]byte("txo"+outTxHash.ToHexString()+fmt.Sprint(outInd))), nil)
	for iter.Next() {
		var utxo TXO
		isConv := byteArr.FromByteArr(iter.Value(), &utxo)
		if !isConv {
			return txoRes, false
		}
		txoRes = append(txoRes, utxo)
	}

	iter.Release()
	err := iter.Error()
	if err != nil {
		return txoRes, false
	}

	return txoRes, true
}

func GetTxoList(prefix string) ([]TXO, bool) {
	var txoArr []TXO
	var isConv bool
	iter := dbController.DB.Db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	for iter.Next() {
		var currTxo TXO
		isConv = byteArr.FromByteArr(iter.Value(), &currTxo)
		if !isConv {
			return txoArr, false
		}
		txoArr = append(txoArr, currTxo)
	}
	iter.Release()
	return txoArr, true
}

func RemoveTxo(outTxHash byteArr.ByteArr, outInd int, BlockHeight uint64) bool {
	return dbController.DB.Db.Delete([]byte("txo"+
		outTxHash.ToHexString()+fmt.Sprint(outInd)+fmt.Sprint(BlockHeight)), nil) == nil
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
