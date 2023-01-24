package test

/*
	STEP 1
	Make emitation of full node work: fill it with accounts, create random transactions
	(valid/not and miscellaneous parameters), create and validate blocks, choose best-fee
	transactions for blocks, change difficulty, log all the information into console,
	control the execution of test by typing in commands.
*/

import (
	"bavovnacoin/account"
	"bavovnacoin/blockchain"
	"bavovnacoin/hashing"
	"bavovnacoin/transaction"
	"fmt"
	"math/rand"
	"time"
)

func initAccountData(leftBound int, rightBound int) {
	source := rand.NewSource(time.Now().Unix())
	rand := rand.New(source)

	accNum := rand.Intn(rightBound-leftBound) + leftBound

	for i := 0; i < accNum; i++ {
		newAcc := account.GenAccount(fmt.Sprint(len(network_accounts)))
		network_accounts = append(network_accounts, newAcc)
	}
}

func getTxRandOuts(currInd int, balance uint64) ([]string, []uint64) {
	source := rand.NewSource(time.Now().Unix())
	rand := rand.New(source)

	accNum := rand.Intn(len(network_accounts)-1) + 1
	var outputAddress []string
	var outputSum []uint64

	for len(outputAddress) < accNum {
		netwAccInd := rand.Intn(len(network_accounts) - 1)
		netwAccAddrInd := rand.Intn(len(network_accounts[netwAccInd].KeyPairList))
		outAddress := hashing.SHA1(network_accounts[netwAccInd].KeyPairList[netwAccAddrInd].PublKey)

		// Check if the same address is already in output of the same tx
		for i := 0; i < len(outputAddress); i++ {
			if outAddress == outputAddress[i] { //TODO: check is user sends his money to his address
				continue
			}
		}
		var output transaction.Output
		output.HashAdr = outAddress

		// Allow spend all money on one tx
		isAllBalanceOnTx := rand.Intn(4)
		if isAllBalanceOnTx == 0 && balance%uint64(accNum) != 0 {
			if len(outputSum)+1 != accNum { // is last output
				output.Sum = uint64(balance / uint64(accNum))
			} else {
				output.Sum = balance - uint64(balance/uint64(accNum))*uint64(accNum)
			}
		} else {
			output.Sum = uint64(float64(balance) / float64(accNum+2))
		}
		outputAddress = append(outputAddress, output.HashAdr)
		outputSum = append(outputSum, output.Sum)
	}
	return outputAddress, outputSum
}

func createRandomTransactions() ([]transaction.Transaction, []bool) {
	var txCorrectness []bool
	var transactions []transaction.Transaction
	for i := 0; i < len(network_accounts); i++ {
		account.CurrAccount = network_accounts[i]
		account.GetBalance()

		if account.CurrAccount.Balance != 0 {
			fee := rand.Intn(5) + 1

			isGenLocktime := rand.Intn(4)
			var locktime uint
			if isGenLocktime == 0 {
				locktime = uint(len(blockchain.Blockchain) + rand.Intn(3) + 1)
			}

			outAddr, outSum := getTxRandOuts(i, account.CurrAccount.Balance)
			tx, mes := transaction.CreateTransaction(fmt.Sprint(i), outAddr, outSum, fee, locktime)

			// Creation od invalid transaction
			isTxInvalid := rand.Intn(5)
			if isTxInvalid == 0 {
				tx.Outputs[0].Sum = account.CurrAccount.Balance
			}

			transactions = append(transactions, tx)
			if len(mes) == 0 && isTxInvalid != 0 {
				txCorrectness = append(txCorrectness, true)
			} else {
				txCorrectness = append(txCorrectness, false)
			}

			network_accounts[i] = account.CurrAccount
		}
	}
	return transactions, txCorrectness
}
