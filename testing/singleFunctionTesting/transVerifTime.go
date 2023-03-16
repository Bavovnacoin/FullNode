package singleFunctionTesting

import (
	"bavovnacoin/account"
	"bavovnacoin/dbController"
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

/*
Test creates a predefined ammount of transactions with a predefined
ammount of incorrect transactions and vereficated them.

The result is printed to the console.
*/
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
	transaction.PrintTransaction(tx[0])
	//tvt.printResults()
	dbController.DB.CloseDb()
	os.RemoveAll(dbController.DbPath)
}
