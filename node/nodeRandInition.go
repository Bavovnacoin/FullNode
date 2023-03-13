package node

import (
	"bavovnacoin/account"
	"bavovnacoin/blockchain"
	"bavovnacoin/byteArr"
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
		newAcc := account.GenAccount(fmt.Sprint(len(account.Wallet)))
		account.Wallet = append(account.Wallet, newAcc)

		account.WriteAccounts()

		newAccountNotFact = rand.Intn(2000000) + 3000000
		newAccountNotCome = 0

		log.Printf("Account with index %d created!\n", len(account.Wallet)-1)
	}
	command_executor.PauseCommand()
}

func getTxRandOuts(currInd int, balance uint64) ([]byteArr.ByteArr, []uint64) {
	source := rand.NewSource(time.Now().Unix())
	rand := rand.New(source)

	var accNum int
	if len(account.Wallet) == 1 {
		accNum = rand.Intn(len(account.Wallet)) + 1
	} else {
		accNum = rand.Intn(len(account.Wallet)-1) + 1
	}
	var outputAddress []byteArr.ByteArr
	var outputValue []uint64

	for len(outputAddress) < accNum {
		netwAccInd := rand.Intn(len(account.Wallet))
		netwAccAddrInd := rand.Intn(len(account.Wallet[netwAccInd].KeyPairList))

		var outAddress byteArr.ByteArr
		outAddress.SetFromHexString(hashing.SHA1(account.Wallet[netwAccInd].KeyPairList[netwAccAddrInd].PublKey), 20)

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
			if len(outputValue)+1 != accNum { // is last output
				output.Value = uint64(balance / uint64(accNum))
			} else {
				output.Value = balance - uint64(balance/uint64(accNum))*uint64(accNum)
			}
		} else {
			output.Value = uint64(float64(balance) / float64(accNum+2))
		}
		outputAddress = append(outputAddress, output.Address)
		outputValue = append(outputValue, output.Value)
	}
	return outputAddress, outputValue
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

		accInd := rand.Int() % len(account.Wallet)
		netwAccAddrInd := rand.Intn(len(account.Wallet[accInd].KeyPairList))

		var accAddr byteArr.ByteArr
		accAddr.SetFromHexString(hashing.SHA1(account.Wallet[accInd].KeyPairList[netwAccAddrInd].PublKey), 20)
		isAddrInMempool := blockchain.IsAddressInMempool(accAddr)

		account.CurrAccount = account.Wallet[accInd]
		account.GetBalance()

		if account.CurrAccount.Balance != 0 && !isAddrInMempool {
			fee := rand.Intn(5) + 1
			isGenLocktime := rand.Intn(5)
			var locktime uint
			if isGenLocktime == 2 {
				locktime = uint(int(blockchain.BcLength+1) + rand.Intn(3) + 1)
			}

			outAddr, outValue := getTxRandOuts(accInd, account.CurrAccount.Balance)
			tx, mes := transaction.CreateTransaction(fmt.Sprint(accInd), outAddr, outValue, fee, locktime)

			// Creation of invalid transaction
			isTxInvalid := rand.Intn(10)
			if isTxInvalid == 1 {
				tx.Outputs[0].Value = account.CurrAccount.Balance
			}

			if len(mes) != 0 || isTxInvalid == 1 {
				log.Println("Created incorrect tx")
			}

			command_executor.PauseCommand()

			account.Wallet[accInd] = account.CurrAccount
			if tx.Inputs != nil {
				if blockchain.AddTxToMempool(tx, true) {
					log.Println("New tx added to mempool")
				} else if isTxInvalid == 1 {
					log.Println("New tx was not added to mempool")
				}
			}
		}
		command_executor.PauseCommand()
	}
}

var AllowCreateBlock bool = true
var CreatedBlock blockchain.Block

func AddBlock(allowLogPrint bool) bool {
	if AllowCreateBlock {
		if allowLogPrint {
			log.Println("Creating a new block")
		}
		go CreateBlockLog(blockchain.GetBits(allowLogPrint), allowLogPrint)
		AllowCreateBlock = false
	}

	if CreatedBlock.MerkleRoot != "" { // Is block mined check
		isBlockValid := blockchain.ValidateBlock(CreatedBlock, int(blockchain.BcLength), true, false)
		AddBlockLog(allowLogPrint, isBlockValid)
		CreatedBlock.MerkleRoot = ""
		return true
	}
	command_executor.PauseCommand()
	return false
}

func CreateBlockLog(bits uint64, allowPrint bool) {
	var rewardAdr byteArr.ByteArr
	rewardAdr.SetFromHexString(blockchain.RewardAddress, 20)
	newBlock := blockchain.CreateBlock(rewardAdr, allowPrint)
	newBlock.Bits = bits
	newBlock = blockchain.MineBlock(newBlock, 1, allowPrint)
	CreatedBlock = newBlock
	command_executor.PauseCommand()
}

func AddBlockLog(allowPrint bool, isBlockValid bool) bool {
	isBlockAdded := false

	if isBlockValid {
		blockchain.AddBlockToBlockchain(CreatedBlock, true)
		if allowPrint {
			log.Println("Block is added to blockchain. Current height: " + fmt.Sprint(blockchain.BcLength+1) + "\n")
		}
		blockchain.IncrBcHeight()
		isBlockAdded = true
	} else {
		if allowPrint {
			log.Println("Block is not added\n")
		}
		isBlockAdded = false
	}

	AllowCreateBlock = true
	CreatedBlock.MerkleRoot = ""
	command_executor.PauseCommand()
	return isBlockAdded
}
