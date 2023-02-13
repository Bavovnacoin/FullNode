package testing

import (
	"bavovnacoin/account"
	"bavovnacoin/blockchain"
	"bavovnacoin/byteArr"
	"bavovnacoin/ecdsa"
	"bavovnacoin/hashing"
	"bavovnacoin/node_controller/command_executor"
	"bavovnacoin/transaction"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var newAccountNotFact int = 0
var newAccountNotCome int = 0

func InitAccountsData() {
	var emptyDat []byte
	err := account.IsWalletExists(account.WalletName)
	if !err {
		os.WriteFile(account.WalletName, emptyDat, 0777)
		var genesisAccKeyPair []ecdsa.KeyPair
		genesisAccKeyPair = append(genesisAccKeyPair, ecdsa.KeyPair{PrivKey: "d966fded26f23d50bb1223cdc6efe4cfebc9f2d6967cb570122c040baf5d42091953a2ba6466963351a4c6bc616e1858de87de02724cc89d9306a62b6d29fab6",
			PublKey: "033361587c679cf9476949cb7cdd15c07d6f2f9674886333f667bfedb87635a4b4"})
		account.Wallet = append(account.Wallet, account.Account{Id: "0",
			HashPass: "b6589fc6ab0dc82cf12099d1c2d40ab994e8410c", KeyPairList: genesisAccKeyPair})

		account.WriteAccounts()

		log.Println("Created new wallet")
	} else {
		dataByte, _ := os.ReadFile(account.WalletName)
		json.Unmarshal([]byte(dataByte), &account.Wallet)
		account.RightBoundAccNum, _ = strconv.Atoi(account.Wallet[len(account.Wallet)-1].Id)

		log.Println("Restored wallet")
	}
}

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
	var outputSum []uint64

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

			outAddr, outSum := getTxRandOuts(accInd, account.CurrAccount.Balance)
			tx, mes := transaction.CreateTransaction(fmt.Sprint(accInd), outAddr, outSum, fee, locktime)

			// Creation of invalid transaction
			isTxInvalid := rand.Intn(10)
			if isTxInvalid == 1 {
				tx.Outputs[0].Sum = account.CurrAccount.Balance
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
	var rewardAdr byteArr.ByteArr
	rewardAdr.SetFromHexString("e930fca003a4a70222d916a74cc851c3b3a9b050", 20)
	createdBlock = blockchain.CreateBlock(rewardAdr, 1)
	command_executor.PauseCommand()
}

func addBlockLog() {
	if blockchain.AddBlockToBlockchain(createdBlock) {
		log.Println("Block is added to blockchain. Current height: " + fmt.Sprint(blockchain.BcLength+1) + "\n")
	} else {
		log.Println("Block is not added\n")
	}
	blockchain.IncrBcHeight()

	allowCreateBlock = true
	createdBlock.MerkleRoot = ""
	command_executor.PauseCommand()
}
