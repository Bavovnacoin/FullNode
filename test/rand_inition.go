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

// func initAccountData(leftBound int, rightBound int) {
// 	source := rand.NewSource(time.Now().Unix())
// 	rand := rand.New(source)

// 	accNum := rand.Intn(rightBound-leftBound) + leftBound

// 	for i := 0; i < accNum; i++ {
// 		newAcc := account.GenAccount(fmt.Sprint(len(network_accounts)))
// 		network_accounts = append(network_accounts, newAcc)
// 	}
// }

var newAccountNotFact int = 0
var newAccountNotCome int = 0

func createAccoundRandom() {
	source := rand.NewSource(time.Now().Unix())
	rand := rand.New(source)

	if newAccountNotCome != newAccountNotFact {
		newAccountNotCome++
	} else {
		newAcc := account.GenAccount(fmt.Sprint(len(network_accounts)))
		network_accounts = append(network_accounts, newAcc)

		newAccountNotFact = rand.Intn(60000) + 100000
		newAccountNotCome = 0
		println("Account created! ")
	}
}

func getTxRandOuts(currInd int, balance uint64) ([]string, []uint64) {
	source := rand.NewSource(time.Now().Unix())
	rand := rand.New(source)

	var accNum int
	if len(network_accounts) == 1 {
		accNum = rand.Intn(len(network_accounts)) + 1
	} else {
		accNum = rand.Intn(len(network_accounts)-1) + 1
	}
	var outputAddress []string
	var outputSum []uint64

	for len(outputAddress) < accNum {
		netwAccInd := rand.Intn(len(network_accounts))
		netwAccAddrInd := rand.Intn(len(network_accounts[netwAccInd].KeyPairList))
		outAddress := hashing.SHA1(network_accounts[netwAccInd].KeyPairList[netwAccAddrInd].PublKey)

		if currInd == netwAccInd {
			continue
		}

		// Check if the same address is already in output of the same tx
		for i := 0; i < len(outputAddress); i++ {
			if outAddress == outputAddress[i] {
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

var newTransNotFact int = 0
var newTransNotCome int = 0

func createRandomTransaction() (transaction.Transaction, bool) {
	source := rand.NewSource(time.Now().Unix())
	rand := rand.New(source)
	var txCorrectness bool
	var newTx transaction.Transaction

	if newTransNotCome != newTransNotFact {
		newTransNotCome++
	} else {
		accInd := rand.Int() % len(network_accounts)
		netwAccAddrInd := rand.Intn(len(network_accounts[accInd].KeyPairList))
		accAddr := hashing.SHA1(network_accounts[accInd].KeyPairList[netwAccAddrInd].PublKey)
		isAddrInMempool := blockchain.IsAddressInMempool(accAddr)

		if isAddrInMempool {
			return newTx, txCorrectness
		}

		account.CurrAccount = network_accounts[accInd]
		account.GetBalance()

		if account.CurrAccount.Balance != 0 {
			fee := rand.Intn(5) + 1
			isGenLocktime := rand.Intn(5)
			var locktime uint
			if isGenLocktime == 1 {
				locktime = uint(len(blockchain.Blockchain) + rand.Intn(3) + 1)
			}

			outAddr, outSum := getTxRandOuts(accInd, account.CurrAccount.Balance)
			tx, mes := transaction.CreateTransaction(fmt.Sprint(accInd), outAddr, outSum, fee, locktime)

			// Creation of invalid transaction
			isTxInvalid := rand.Intn(5)
			if isTxInvalid == 1 {
				tx.Outputs[0].Sum = account.CurrAccount.Balance
			}

			newTx = tx
			if len(mes) == 0 && isTxInvalid != 1 {
				txCorrectness = true
			} else {
				txCorrectness = false
			}

			network_accounts[accInd] = account.CurrAccount

			newTransNotFact = rand.Intn(30000) + 10000
			newTransNotCome = 0
			//println(accAddr)
		}
	}
	return newTx, txCorrectness
}

var allowCreateBlock bool = true
var createdBlock blockchain.Block

func addBlock() {
	if allowCreateBlock {
		go createBlockLog()
		allowCreateBlock = false
	}
	if createdBlock.MerkleRoot != "" {
		addBlockLog()
	}
}

func createBlockLog() {
	println("Creating a block with mempool len " + fmt.Sprint(len(blockchain.Mempool)))
	createdBlock = blockchain.CreateBlock(len(blockchain.Blockchain), "e930fca003a4a70222d916a74cc851c3b3a9b050", 1)
	//allowCreateBlock = false
}

func addBlockLog() {
	if blockchain.AddBlockToBlockchain(createdBlock) {
		println("Block is added to blockchain")
	} else {
		println("Block is not added")
	}
	println(fmt.Sprint(len(blockchain.Blockchain)) + " - blockchain length")
	allowCreateBlock = true
	createdBlock.MerkleRoot = ""
}
