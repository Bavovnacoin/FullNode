package testing

/*
	STEP 1
	Make emitation of full node work: fill it with accounts, create random transactions
	(valid/not and miscellaneous parameters), create and validate blocks, choose best-fee
	transactions for blocks, change difficulty, log all the information into console,
	control the execution of test by typing in commands.
*/

import (
	"bavovnacoin/account"
	"bavovnacoin/address"
	"bavovnacoin/blockchain"
	"bavovnacoin/hashing"
	"bavovnacoin/node_controller/command_executor"
	"bavovnacoin/transaction"
	"fmt"
	"log"
	"math/rand"
	"time"
)

var newAccountNotFact int = 0
var newAccountNotCome int = 0

func createAccoundRandom() {
	source := rand.NewSource(time.Now().Unix())
	rand := rand.New(source)

	if newAccountNotCome != newAccountNotFact {
		newAccountNotCome++
	} else {
		newAcc := account.GenAccount(fmt.Sprint(len(command_executor.Network_accounts)))
		command_executor.Network_accounts = append(command_executor.Network_accounts, newAcc)

		newAccountNotFact = rand.Intn(2000000) + 3000000
		newAccountNotCome = 0
		log.Printf("Account with index %d created!\n", len(command_executor.Network_accounts)-1)
	}
	command_executor.PauseCommand()
}

func getTxRandOuts(currInd int, balance uint64) ([]address.Address, []uint64) {
	source := rand.NewSource(time.Now().Unix())
	rand := rand.New(source)

	var accNum int
	if len(command_executor.Network_accounts) == 1 {
		accNum = rand.Intn(len(command_executor.Network_accounts)) + 1
	} else {
		accNum = rand.Intn(len(command_executor.Network_accounts)-1) + 1
	}
	var outputAddress []address.Address
	var outputSum []uint64

	for len(outputAddress) < accNum {
		netwAccInd := rand.Intn(len(command_executor.Network_accounts))
		netwAccAddrInd := rand.Intn(len(command_executor.Network_accounts[netwAccInd].KeyPairList))

		var outAddress address.Address
		outAddress.SetFromHexString(hashing.SHA1(command_executor.Network_accounts[netwAccInd].KeyPairList[netwAccAddrInd].PublKey))

		if currInd == netwAccInd {
			continue
		}

		// Check if the same address is already in output of the same tx
		for i := 0; i < len(outputAddress); i++ {
			if outAddress.IsEqual(outputAddress[i]) {
				continue
			}
		}
		var output transaction.Output
		output.Address = outAddress

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
		outputAddress = append(outputAddress, output.Address)
		outputSum = append(outputSum, output.Sum)
	}
	return outputAddress, outputSum
}

func txRandomCreator() {
	go createTxRandom()
}

var sleepTimeTxCreation uint64 = 10

func createTxRandom() {
	for command_executor.Node_working {
		source := rand.NewSource(time.Now().Unix())
		rand := rand.New(source)

		time.Sleep(time.Duration(sleepTimeTxCreation) * time.Millisecond)
		sleepTimeTxCreation = uint64(rand.Intn(300)) + 1000

		accInd := rand.Int() % len(command_executor.Network_accounts)
		netwAccAddrInd := rand.Intn(len(command_executor.Network_accounts[accInd].KeyPairList))

		var accAddr address.Address
		accAddr.SetFromHexString(hashing.SHA1(command_executor.Network_accounts[accInd].KeyPairList[netwAccAddrInd].PublKey))
		isAddrInMempool := blockchain.IsAddressInMempool(accAddr)

		account.CurrAccount = command_executor.Network_accounts[accInd]
		account.GetBalance()

		if account.CurrAccount.Balance != 0 && !isAddrInMempool {
			fee := rand.Intn(5) + 1
			isGenLocktime := rand.Intn(5)
			var locktime uint
			if isGenLocktime == 2 {
				locktime = uint(len(blockchain.Blockchain) + rand.Intn(3) + 1)
			}

			outAddr, outSum := getTxRandOuts(accInd, account.CurrAccount.Balance)
			tx, mes := transaction.CreateTransaction(fmt.Sprint(accInd), outAddr, outSum, fee, locktime)

			// Creation of invalid transaction
			isTxInvalid := rand.Intn(5)
			if isTxInvalid == 1 {
				tx.Outputs[0].Sum = account.CurrAccount.Balance
			}

			if len(mes) != 0 || isTxInvalid == 1 {
				log.Println("Created incorrect tx")
			}

			command_executor.PauseCommand()

			command_executor.Network_accounts[accInd] = account.CurrAccount
			if tx.Inputs != nil {
				if blockchain.AddTxToMempool(tx) {
					log.Println("New tx added to mempool")
				} else if isTxInvalid == 1 {
					log.Println("New tx was not added to mempool")
				}
			}
		}
		command_executor.PauseCommand()
	}
}

var allowCreateBlock bool = true
var createdBlock blockchain.Block

func addBlock() {
	if allowCreateBlock {
		log.Println("Creating a new block")
		go createBlockLog()
		allowCreateBlock = false
	}
	if createdBlock.MerkleRoot != "" {
		addBlockLog()
	}
	command_executor.PauseCommand()
}

func createBlockLog() {
	var rewardAdr address.Address
	rewardAdr.SetFromHexString("e930fca003a4a70222d916a74cc851c3b3a9b050")
	createdBlock = blockchain.CreateBlock(rewardAdr, 1)
	command_executor.PauseCommand()
}

func addBlockLog() {
	if blockchain.AddBlockToBlockchain(createdBlock) {
		log.Println("Block is added to blockchain. Current length: " + fmt.Sprint(len(blockchain.Blockchain)) + "\n")
	} else {
		log.Println("Block is not added\n")
	}

	allowCreateBlock = true
	createdBlock.MerkleRoot = ""
	command_executor.PauseCommand()
}
