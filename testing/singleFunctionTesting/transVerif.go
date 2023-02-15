package singleFunctionTesting

import (
	"bavovnacoin/account"
	"bavovnacoin/byteArr"
	"bavovnacoin/cryption"
	"bavovnacoin/ecdsa"
	"bavovnacoin/hashing"
	"bavovnacoin/transaction"
	"bavovnacoin/txo"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

var txAmmount int = 10         // Total ammount of txs
var incorrectTxAmmount int = 3 // Total ammount of incorrect txs

var testTransactions []transaction.Transaction
var txIncorrMessage []string

var source rand.Source
var random *rand.Rand

func genTestAccounts() {
	for i := 0; i < txAmmount; i++ {
		account.Wallet = append(account.Wallet, account.GenAccount(fmt.Sprint(i)))
	}
}

func genTestUtxo() {
	for i := 0; i < txAmmount; i++ {
		var outAddr byteArr.ByteArr
		outAddr.SetFromHexString(hashing.SHA1(account.Wallet[i].KeyPairList[0].PublKey), 20)

		var outTxHash byteArr.ByteArr
		outTxHash.SetFromHexString(hashing.SHA1(fmt.Sprint(i)+fmt.Sprint(time.Now().Unix())), 20)

		newTxo := txo.TXO{OutTxHash: outTxHash, Value: (rand.Uint64()%700000 + 300000), OutAddress: outAddr}
		txo.CoinDatabase = append(txo.CoinDatabase, newTxo)
	}
}

func makeTxIncorrect(tx transaction.Transaction, incTxCounter int) (transaction.Transaction, string) {
	if incTxCounter%3 == 0 { // Wrong tx hash in input
		var incTxHash byteArr.ByteArr
		inpId := 0
		if len(tx.Inputs) != 0 {
			inpId = rand.Intn(len(tx.Inputs))
		}
		tx.Inputs[inpId].TxHash = incTxHash

		return tx, "Wrong tx hash in input num. " + fmt.Sprint(inpId)
	} else if incTxCounter%3 == 1 { // Wrong tx signature in input
		var incTxSig byteArr.ScriptSig
		fakeSigAccId := 0
		inpId := 0
		if len(tx.Inputs) != 0 {
			inpId = rand.Intn(len(tx.Inputs))
			fakeSigAccId = rand.Intn(len(tx.Inputs))
		}

		fakeSigPrivKey := cryption.AES_decrypt(account.Wallet[fakeSigAccId].KeyPairList[0].PrivKey, fmt.Sprint(fakeSigAccId))
		fakeSig := ecdsa.Sign(hashing.SHA1("Glory to Ukraine"), fakeSigPrivKey)

		incTxSig.SetFromHexString(tx.Inputs[inpId].ScriptSig.GetPubKey().ToHexString()+fakeSig, 111)
		tx.Inputs[inpId].ScriptSig = incTxSig
		return tx, "Wrong signature in input num. " + fmt.Sprint(inpId)
	} else { // Wrong tx output value
		outId := 0
		if len(tx.Inputs) != 0 {
			outId = rand.Intn(len(tx.Outputs))
		}

		tx.Outputs[outId].Value = ^uint64(0)
		return tx, "Too high output value in output num. " + fmt.Sprint(outId)
	}
}

func genRandTxs() {
	var endStep int = txAmmount / incorrectTxAmmount
	var incTxInd int = rand.Intn(endStep)
	var incTxCounter int

	for i := 0; i < txAmmount; i++ {
		account.CurrAccount = account.Wallet[i]
		account.GetBalance()
		var outAddr byteArr.ByteArr
		outAddr.SetFromHexString(hashing.SHA1(account.Wallet[i].KeyPairList[0].PublKey), 20)

		var outAddrTx []byteArr.ByteArr
		outAddrTx = append(outAddrTx, outAddr)

		var outValTx []uint64
		outValTx = append(outValTx, uint64(txo.CoinDatabase[i].Value/((rand.Uint64()%10)+3)))

		newTx, _ := transaction.CreateTransaction(fmt.Sprint(i), outAddrTx, outValTx, rand.Intn(10), uint(rand.Intn(10)))

		if i == incTxInd {
			endStep := i + int(txAmmount/incorrectTxAmmount)
			r := rand.Intn(endStep)
			incTxInd = r + i + 1
			println(r)
			var message string
			newTx, message = makeTxIncorrect(newTx, incTxCounter)

			txIncorrMessage = append(txIncorrMessage, message)
			incTxCounter++
		} else {
			txIncorrMessage = append(txIncorrMessage, "")
		}

		testTransactions = append(testTransactions, newTx)
	}
}

func incorMessageToBool(message string) bool {
	if message == "" {
		return true
	}
	return false
}

func printResults() {
	log.Printf("Tx index %s Problem %s Verif. result %s Is matched\n", strings.Repeat(" ", 6), strings.Repeat(" ", 36), strings.Repeat(" ", 6))
	resultNotMatchedCounter := 0
	for i := 0; i < txAmmount; i++ {
		verifMes := "Correct"
		verifRes := transaction.VerifyTransaction(testTransactions[i])
		if !verifRes {
			verifMes = "Incorrect"
		}

		incorMes := "-"
		if txIncorrMessage[i] != "" {
			incorMes = txIncorrMessage[i]
		}

		isMatchedMes := "Yes"
		if incorMessageToBool(txIncorrMessage[i]) != verifRes {
			isMatchedMes = "No"
			resultNotMatchedCounter++
		}
		log.Printf("   [%s %s %s %s", fmt.Sprint(i)+"]"+strings.Repeat(" ", 10-len(fmt.Sprint(i))), incorMes+strings.Repeat(" ", 44-len(incorMes)),
			verifMes+strings.Repeat(" ", 20-len(verifMes)), isMatchedMes)
	}

	log.Printf("Test result: %d\\%d\n", txAmmount-resultNotMatchedCounter, txAmmount)
}

/*
Test creates a predefined ammount of transactions with a predefined
ammount of incorrect transactions and vereficated them.

The result is printed to the console.
*/
func TransactionsVerefication() {
	source = rand.NewSource(time.Now().Unix())
	random = rand.New(source)

	if incorrectTxAmmount > txAmmount {
		log.Println("Wrong input")
	}

	genTestAccounts()
	log.Printf("Generated %d test accounts\n", len(account.Wallet))
	genTestUtxo()
	log.Printf("Generated %d test utxo\n", len(txo.CoinDatabase))
	genRandTxs()
	printResults()
}
