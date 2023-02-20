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

type TransVerifTest struct {
	txAmmount          int // Total ammount of txs
	incorrectTxAmmount int // Total ammount of incorrect txs

	testTransactions []transaction.Transaction
	txIncorrMessage  []string

	source rand.Source
	random *rand.Rand
}

func (tvt *TransVerifTest) SetTestValues(txAmmount int, incorrectTxAmmount int) {
	tvt.txAmmount = txAmmount
	tvt.incorrectTxAmmount = incorrectTxAmmount
}

func (tvt *TransVerifTest) genTestAccounts(ammount int) {
	for i := 0; i < ammount; i++ {
		account.Wallet = append(account.Wallet, account.GenAccount(fmt.Sprint(i)))
	}
}

func (tvt *TransVerifTest) genTestUtxo(ammount int) {
	for i := 0; i < ammount; i++ {
		var outAddr byteArr.ByteArr
		outAddr.SetFromHexString(hashing.SHA1(account.Wallet[i].KeyPairList[0].PublKey), 20)

		var outTxHash byteArr.ByteArr
		outTxHash.SetFromHexString(hashing.SHA1(fmt.Sprint(i)+fmt.Sprint(time.Now().Unix())), 20)

		newTxo := txo.TXO{OutTxHash: outTxHash, Value: (tvt.random.Uint64()%700000 + 300000), OutAddress: outAddr}
		txo.CoinDatabase = append(txo.CoinDatabase, newTxo)
	}
}

func (tvt *TransVerifTest) makeTxIncorrect(tx transaction.Transaction, incTxCounter int) (transaction.Transaction, string) {
	if incTxCounter%3 == 0 { // Wrong tx hash in input
		var incTxHash byteArr.ByteArr
		inpId := 0
		if len(tx.Inputs) != 0 {
			inpId = tvt.random.Intn(len(tx.Inputs))
		}
		tx.Inputs[inpId].TxHash = incTxHash

		return tx, "Wrong tx hash in input num. " + fmt.Sprint(inpId)
	} else if incTxCounter%3 == 1 { // Wrong tx signature in input
		var incTxSig byteArr.ScriptSig
		fakeSigAccId := 0
		inpId := 0
		if len(tx.Inputs) != 0 {
			inpId = tvt.random.Intn(len(tx.Inputs))
			fakeSigAccId = tvt.random.Intn(len(tx.Inputs))
		}

		fakeSigPrivKey := cryption.AES_decrypt(account.Wallet[fakeSigAccId].KeyPairList[0].PrivKey, fmt.Sprint(fakeSigAccId))
		fakeSig := ecdsa.Sign(hashing.SHA1("Glory to Ukraine"), fakeSigPrivKey)

		incTxSig.SetFromHexString(tx.Inputs[inpId].ScriptSig.GetPubKey().ToHexString()+fakeSig, 111)
		tx.Inputs[inpId].ScriptSig = incTxSig
		return tx, "Wrong signature in input num. " + fmt.Sprint(inpId)
	} else { // Wrong tx output value
		outId := 0
		if len(tx.Inputs) != 0 {
			outId = tvt.random.Intn(len(tx.Outputs))
		}

		tx.Outputs[outId].Value = ^uint64(0)
		return tx, "Too high output value in output num. " + fmt.Sprint(outId)
	}
}

func (tvt *TransVerifTest) genValidTx(currAccId int) transaction.Transaction {
	account.CurrAccount = account.Wallet[currAccId]
	account.GetBalance()
	var outAddr byteArr.ByteArr
	outAddr.SetFromHexString(hashing.SHA1(account.Wallet[currAccId].KeyPairList[0].PublKey), 20)

	var outAddrTx []byteArr.ByteArr
	outAddrTx = append(outAddrTx, outAddr)

	var outValTx []uint64
	outValTx = append(outValTx, uint64(txo.CoinDatabase[currAccId].Value/((tvt.random.Uint64()%10)+3)))

	newTx, _ := transaction.CreateTransaction(fmt.Sprint(currAccId), outAddrTx, outValTx, tvt.random.Intn(10), uint(tvt.random.Intn(10)))
	return newTx
}

func (tvt *TransVerifTest) genRandTxs() {
	var step int = int(tvt.txAmmount / tvt.incorrectTxAmmount)
	var incTxInd int = -1
	var incTxCounter int

	if tvt.incorrectTxAmmount != 0 {
		stStep := step * incTxCounter
		incTxInd = tvt.random.Intn(step) + stStep

		incTxCounter++
	}

	for i := 0; i < tvt.txAmmount; i++ {
		newTx := tvt.genValidTx(i)

		if i == incTxInd && incTxCounter <= tvt.incorrectTxAmmount {
			stStep := step * incTxCounter
			if tvt.incorrectTxAmmount-1 == incTxCounter {
				step = tvt.txAmmount - stStep
			}
			incTxInd = tvt.random.Intn(step) + stStep

			var message string
			newTx, message = tvt.makeTxIncorrect(newTx, incTxCounter)
			tvt.txIncorrMessage = append(tvt.txIncorrMessage, message)
			incTxCounter++
		} else {
			tvt.txIncorrMessage = append(tvt.txIncorrMessage, "")
		}

		tvt.testTransactions = append(tvt.testTransactions, newTx)
	}
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
func (tvt *TransVerifTest) TransactionsVerefication() {
	tvt.source = rand.NewSource(time.Now().Unix())
	tvt.random = rand.New(tvt.source)

	if tvt.incorrectTxAmmount > tvt.txAmmount {
		log.Println("Wrong input")
	}

	tvt.genTestAccounts(tvt.txAmmount)
	log.Printf("Generated %d test accounts\n", len(account.Wallet))
	tvt.genTestUtxo(tvt.txAmmount)
	log.Printf("Generated %d test utxo\n", len(txo.CoinDatabase))
	log.Printf("Generating %d txs (%d are incorrect)", tvt.txAmmount, tvt.incorrectTxAmmount)
	tvt.genRandTxs()
	tvt.printResults()
}
