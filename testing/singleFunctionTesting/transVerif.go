/*
	Checks a node ability to detect incorrect transactions
*/

package singleFunctionTesting

import (
	"bavovnacoin/dbController"
	"bavovnacoin/testing"
	"bavovnacoin/testing/account"
	"bavovnacoin/transaction"
	"bavovnacoin/txo"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type TransVerifTest struct {
	txAmmount          int // Total ammount of txs
	incorrectTxAmmount int // Total ammount of incorrect txs

	testTransactions []transaction.Transaction
	txIncorrMessage  []string

	source rand.Source
	random *rand.Rand
}

func incorMessageToBool(message string) bool {
	if message == "" {
		return true
	}
	return false
}

func (tvt *TransVerifTest) printResults() {
	log.Printf("Tx index %s Problem %s Verif. result %s Is matched\n", strings.Repeat(" ", 6), strings.Repeat(" ", 36), strings.Repeat(" ", 6))
	resultNotMatchedCounter := 0
	for i := 0; i < tvt.txAmmount; i++ {
		verifMes := "Correct"
		verifRes := transaction.VerifyTransaction(tvt.testTransactions[i])
		if !verifRes {
			verifMes = "Incorrect"
		}

		incorMes := "-"
		if tvt.txIncorrMessage[i] != "" {
			incorMes = tvt.txIncorrMessage[i]
		}

		isMatchedMes := "Yes"
		if incorMessageToBool(tvt.txIncorrMessage[i]) != verifRes {
			isMatchedMes = "No"
			resultNotMatchedCounter++
		}
		log.Printf("   [%s %s %s %s", fmt.Sprint(i)+"]"+strings.Repeat(" ", 10-len(fmt.Sprint(i))), incorMes+strings.Repeat(" ", 44-len(incorMes)),
			verifMes+strings.Repeat(" ", 20-len(verifMes)), isMatchedMes)
	}

	result := "Passed"
	if tvt.txAmmount-resultNotMatchedCounter != tvt.txAmmount {
		result = "Not passed"
	}

	log.Printf("Test result: %d\\%d. %s\n", tvt.txAmmount-resultNotMatchedCounter, tvt.txAmmount, result)
}

/*
Test creates a predefined ammount of transactions with a predefined
ammount of incorrect transactions and vereficated them.

The result is printed to the console.
*/
func (tvt *TransVerifTest) TransactionsVerefication(txAmmount int, incorrectTxAmmount int) {
	tvt.txAmmount = txAmmount
	tvt.incorrectTxAmmount = incorrectTxAmmount

	dbController.DbPath = "testing/testData"
	if _, err := os.Stat(dbController.DbPath); err == nil {
		os.RemoveAll(dbController.DbPath)
		println("Removed test db from a previous test.")
	}
	dbController.DB.OpenDb()

	tvt.source = rand.NewSource(time.Now().Unix())
	tvt.random = rand.New(tvt.source)

	if tvt.incorrectTxAmmount > tvt.txAmmount {
		log.Println("Wrong input")
	}

	testing.GenTestAccounts(tvt.txAmmount)
	log.Printf("Generated %d test accounts\n", len(account.Wallet))
	testing.GenTestUtxo(tvt.txAmmount, tvt.random)
	log.Printf("Generated %d test utxo\n", len(txo.CoinDatabase))
	tvt.testTransactions, tvt.txIncorrMessage = testing.GenRandTxs(tvt.txAmmount, tvt.incorrectTxAmmount, tvt.random)
	log.Printf("Generated %d txs (%d are incorrect)", tvt.txAmmount, tvt.incorrectTxAmmount)
	tvt.printResults()
	dbController.DB.CloseDb()
	os.RemoveAll(dbController.DbPath)
}
