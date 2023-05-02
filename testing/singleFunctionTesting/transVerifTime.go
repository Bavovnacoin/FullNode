package singleFunctionTesting

import (
	"bavovnacoin/dbController"
	"bavovnacoin/ecdsa"
	"bavovnacoin/hashing"
	"bavovnacoin/testing"
	"bavovnacoin/testing/account"
	"bavovnacoin/transaction"
	"bavovnacoin/txo"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

type TransVerifTime struct {
	txVerifTime        [4]time.Duration
	commonTimeConsumed time.Duration

	source rand.Source
	random *rand.Rand
}

// Same function as in transaction package but with time checking
func VerifyTransactionTime(tx transaction.Transaction) (bool, [4]time.Duration) {
	var timeFunConsumed [4]time.Duration
	if tx.Version == 0 {
		ecdsa.InitValues()

		var timeStart time.Time
		var inpValue uint64
		var outValue uint64
		hashMesOfTx := hashing.SHA1(transaction.GetCatTxFields(tx))

		// Checking signatures and unique inputs
		for i := 0; i < len(tx.Inputs); i++ {
			timeStart = time.Now() // ToHexString time check
			if len(tx.Inputs[i].ScriptSig.ToHexString()) < 66 {
				timeFunConsumed[0] += time.Since(timeStart)
				return false, timeFunConsumed
			}
			timeFunConsumed[0] += time.Since(timeStart)

			timeStart = time.Now() // GetPubKey and GetSignature time check
			pubKey := tx.Inputs[i].ScriptSig.GetPubKey().ToHexString()
			sign := tx.Inputs[i].ScriptSig.GetSignature().ToHexString()
			timeFunConsumed[1] += time.Since(timeStart)

			timeStart = time.Now() // Sig verify time check
			if !ecdsa.Verify(pubKey, sign, hashMesOfTx) {
				timeFunConsumed[2] += time.Since(timeStart)
				return false, timeFunConsumed
			}
			timeFunConsumed[2] += time.Since(timeStart)

			timeStart = time.Now() // UTXO from db getting time check
			utxos, res := txo.GetUtxos(tx.Inputs[i].TxHash, tx.Inputs[i].OutInd)
			timeFunConsumed[3] += time.Since(timeStart)
			if !res {
				return false, timeFunConsumed
			}
			curVal := utxos[0].Value
			inpValue += curVal
		}
		timeFunConsumed[0] /= time.Duration(len(tx.Inputs))
		timeFunConsumed[1] /= time.Duration(len(tx.Inputs))
		timeFunConsumed[2] /= time.Duration(len(tx.Inputs))
		timeFunConsumed[3] /= time.Duration(len(tx.Inputs))

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

func (tvt *TransVerifTime) printResults() {
	println("Verififcation results:")
	fmt.Printf("Total time consumed (ms): %d\n", tvt.commonTimeConsumed/time.Millisecond)
	fmt.Printf("ToHexString function time %d (mcs)\n", tvt.txVerifTime[0]/time.Microsecond)
	fmt.Printf("GetPubKey and GetSignature function time %d (mcs)\n", tvt.txVerifTime[1]/time.Microsecond)
	fmt.Printf("Sig verify function time %d (mcs)\n", tvt.txVerifTime[2]/time.Microsecond)
	fmt.Printf("UTXO from db getting function time %d (mcs)\n", tvt.txVerifTime[3]/time.Microsecond)
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

	testing.GenTestAccounts(50)
	log.Printf("Generated %d test accounts\n", len(account.Wallet))
	testing.GenTestUtxo(50, tvt.random)
	log.Printf("Generated %d test utxo\n", len(txo.CoinDatabase))

	tx, _ := testing.GenRandTxs(1, 0, tvt.random)
	var verifRes bool

	commonTimeConsumedStart := time.Now()
	verifRes, tvt.txVerifTime = VerifyTransactionTime(tx[0])
	tvt.commonTimeConsumed = time.Since(commonTimeConsumedStart)

	if verifRes {
		tvt.printResults()
	} else {
		println("Incorrect tx")
	}

	dbController.DB.CloseDb()
	os.RemoveAll(dbController.DbPath)
}
