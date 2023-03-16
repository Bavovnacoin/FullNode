package singleFunctionTesting

import (
	"bavovnacoin/account"
	"bavovnacoin/dbController"
	"bavovnacoin/ecdsa"
	"bavovnacoin/hashing"
	"bavovnacoin/testing"
	"bavovnacoin/transaction"
	"bavovnacoin/txo"
	"log"
	"math/rand"
	"os"
	"time"
)

type TransVerifTime struct {
	source rand.Source
	random *rand.Rand
}

// Same function as in transaction package but with time checking
func VerifyTransactionTime(tx transaction.Transaction) (bool, []time.Duration) {
	var timeFunConsumed []time.Duration
	if tx.Version == 0 {
		ecdsa.InitValues()

		var inpValue uint64
		var outValue uint64
		hashMesOfTx := hashing.SHA1(transaction.GetCatTxFields(tx))

		// Checking signatures and unique inputs
		for i := 0; i < len(tx.Inputs); i++ {
			if len(tx.Inputs[i].ScriptSig.ToHexString()) < 66 {
				return false, timeFunConsumed
			}

			pubKey := tx.Inputs[i].ScriptSig.GetPubKey().ToHexString()
			sign := tx.Inputs[i].ScriptSig.GetSignature().ToHexString()

			if !ecdsa.Verify(pubKey, sign, hashMesOfTx) {
				return false, timeFunConsumed
			}

			utxo, res := txo.GetUtxo(tx.Inputs[i].TxHash, tx.Inputs[i].OutInd)
			if !res {
				return false, timeFunConsumed
			}
			curVal := utxo.Value
			//curVal := account.GetBalHashOutInd(tx.Inputs[i].TxHash, tx.Inputs[i].OutInd)
			inpValue += curVal
		}

		for i := 0; i < len(tx.Outputs); i++ {
			outValue += tx.Outputs[i].Value
		}

		// Checking presence of coins to be spent
		if inpValue < outValue {
			return false, timeFunConsumed
		}
		return true, timeFunConsumed
	}
	return false, timeFunConsumed
}

func (tvt *TransVerifTime) TransVerifTime() {

	dbController.DbPath = "testing/testData"
	if _, err := os.Stat(dbController.DbPath); err == nil {
		os.RemoveAll(dbController.DbPath)
		println("Removed test db from a previous test.")
	}
	dbController.DB.OpenDb()

	tvt.source = rand.NewSource(time.Now().Unix())
	tvt.random = rand.New(tvt.source)

	testing.GenTestAccounts(1)
	log.Printf("Generated %d test accounts\n", len(account.Wallet))
	testing.GenTestUtxo(1, tvt.random)
	log.Printf("Generated %d test utxo\n", len(txo.CoinDatabase))

	tx, _ := testing.GenRandTxs(1, 0, tvt.random)
	transaction.VerifyTransaction(tx[0])
	dbController.DB.CloseDb()
	os.RemoveAll(dbController.DbPath)
}
