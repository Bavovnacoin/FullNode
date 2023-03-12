package singleFunctionTesting

import (
	"bavovnacoin/account"
	"bavovnacoin/blockchain"
	"bavovnacoin/byteArr"
	"bavovnacoin/hashing"
	"bavovnacoin/transaction"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"strings"
	"time"
)

type BlockchainVerifTest struct {
	blockAmmount          int // Total ammount of blocks
	incorrectblockAmmount int // Total ammount of incorrect blocks
	txPerBlockAmmount     int // Ammount of transactions per block
	currAccWalletInd      int
	blockIncorrMessage    []string
	factBlockCorrectness  []bool

	source rand.Source
	random *rand.Rand
}

func (bvt *BlockchainVerifTest) SetTestValues(blockAmmount int, incorrectblockAmmount int, txPerBlockAmmount int) {
	bvt.blockAmmount = blockAmmount
	bvt.incorrectblockAmmount = incorrectblockAmmount
	bvt.txPerBlockAmmount = txPerBlockAmmount
}

func (bvt *BlockchainVerifTest) genBlockTestAccounts(ammount int) {
	for i := 0; i < ammount; i++ {
		account.Wallet = append(account.Wallet, account.GenAccount(fmt.Sprint(len(account.Wallet))))
	}
}

func (bvt *BlockchainVerifTest) createValidTx() transaction.Transaction {
	txCreationTryCounter := 0

	for txCreationTryCounter < bvt.txPerBlockAmmount {
		accInd := bvt.currAccWalletInd
		bvt.currAccWalletInd = (bvt.currAccWalletInd + 1) % bvt.txPerBlockAmmount

		var accAddr byteArr.ByteArr
		accAddr.SetFromHexString(hashing.SHA1(account.Wallet[accInd].KeyPairList[0].PublKey), 20)
		isAddrInMempool := blockchain.IsAddressInMempool(accAddr)

		account.CurrAccount = account.Wallet[accInd]
		account.GetBalance()
		account.Wallet[accInd] = account.CurrAccount

		if isAddrInMempool == true || account.CurrAccount.Balance == 0 {
			txCreationTryCounter++
			continue
		}

		fee := bvt.random.Intn(5) + 1
		isGenLocktime := bvt.random.Intn(5)
		var locktime uint
		if isGenLocktime == 2 {
			locktime = uint(int(blockchain.BcLength+1) + bvt.random.Intn(2) + 1)
		}

		var outAddr []byteArr.ByteArr
		outAddr = append(outAddr, byteArr.ByteArr{})
		outAddr[0].SetFromHexString(hashing.SHA1(account.Wallet[bvt.currAccWalletInd].KeyPairList[0].PublKey), 20)

		var outValue []uint64
		outValue = append(outValue, account.CurrAccount.Balance/2)

		tx, _ := transaction.CreateTransaction(fmt.Sprint(accInd), outAddr, outValue, fee, locktime)

		if tx.Inputs != nil {
			return tx
		}

	}

	var invalidTx transaction.Transaction // Enough ammount of transactions is already in mempool.
	return invalidTx
}

func (bvt *BlockchainVerifTest) addMempValidTxs(currBlockTxAmmount int) {
	for i := 0; i < currBlockTxAmmount; i++ {
		blockchain.AddTxToMempool(bvt.createValidTx(), true)
	}
}

func (bvt *BlockchainVerifTest) makeBlockIncorrect(block blockchain.Block, incBlockCounter int) (blockchain.Block, string) {
	if incBlockCounter%5 == 1 { // Wrong nonce
		block.Nonce = 0
		for true {
			blockHash := hashing.SHA1(blockchain.BlockToString(block))
			bigBlockHash, _ := new(big.Int).SetString(blockHash, 16)
			if blockchain.BitsToTarget(block.Bits).Cmp(bigBlockHash) == -1 {
				return block, "Wrong nonce value"
			}
			block.Nonce++
		}
	} else if incBlockCounter%5 == 2 { // Wrong coinbase tx
		block.Transactions[0].Outputs[0].Value *= 2
		return block, "Wrong coinbase value"
	} else if incBlockCounter%5 == 3 { // Wrong merkle root
		block.MerkleRoot = hashing.SHA1("Glory to Ukraine")
		return block, "Wrong merkle root value"
	} else if incBlockCounter%5 == 4 { // Wrong prev block hash
		block.HashPrevBlock = hashing.SHA1("Glory to Ukraine")
		return block, "Wrong hash of previous block value"
	} else { // Wrong tx
		if len(block.Transactions) == 1 {
			block.Transactions[0].Outputs[0].Value *= 2
			return block, "Wrong coinbase value"
		}
		block.Transactions[len(block.Transactions)-1].Outputs[0].Value *= 100000000
		return block, "Wrong transaction"
	}
	return block, ""
}

func (bvt *BlockchainVerifTest) genBlocks() {
	var step int = int(bvt.blockAmmount / bvt.incorrectblockAmmount)
	var incBlockInd int = -1
	var incBlockCounter int
	accId := 0 // Account id for reward when block created

	if bvt.incorrectblockAmmount != 0 {
		stStep := step * incBlockCounter
		incBlockInd = bvt.random.Intn(step) + stStep
		incBlockCounter++
	}

	for i := 0; i < bvt.blockAmmount; i++ {
		currBlockTxAmmount := bvt.random.Intn(bvt.txPerBlockAmmount)
		bvt.addMempValidTxs(currBlockTxAmmount)

		var rewAddress byteArr.ByteArr
		rewAddress.SetFromHexString(hashing.SHA1(account.Wallet[accId].KeyPairList[0].PublKey), 20)
		block := blockchain.CreateBlock(rewAddress, false)
		block.Bits = 0xf00fff14
		block = blockchain.MineBlock(block, 1, false)

		if i == incBlockInd && incBlockCounter <= bvt.incorrectblockAmmount {
			stStep := step * incBlockCounter
			if bvt.incorrectblockAmmount-1 == incBlockCounter {
				step = bvt.blockAmmount - stStep
			}
			incBlockInd = bvt.random.Intn(step) + stStep
			var message string
			block, message = bvt.makeBlockIncorrect(block, incBlockCounter)
			bvt.blockIncorrMessage = append(bvt.blockIncorrMessage, message)
			incBlockCounter++
		} else {
			bvt.blockIncorrMessage = append(bvt.blockIncorrMessage, "")
		}

		if blockchain.AddBlockToBlockchain(block, false) {
			bvt.factBlockCorrectness = append(bvt.factBlockCorrectness, true)
			blockchain.BcLength++
		} else {
			bvt.factBlockCorrectness = append(bvt.factBlockCorrectness, false)
		}

		accId = (accId + 1) % bvt.txPerBlockAmmount
	}
}

func (bvt *BlockchainVerifTest) printResults() {
	log.Printf("Tx index %s Problem %s Verif. result %s Is matched\n", strings.Repeat(" ", 6), strings.Repeat(" ", 36), strings.Repeat(" ", 6))
	resultNotMatchedCounter := 0
	for i := 0; i < bvt.blockAmmount; i++ {
		verifMes := "Correct"
		if !bvt.factBlockCorrectness[i] {
			verifMes = "Incorrect"
		}

		incorMes := "-"
		if bvt.blockIncorrMessage[i] != "" {
			incorMes = bvt.blockIncorrMessage[i]
		}

		isMatchedMes := "Yes"
		if incorMessageToBool(bvt.blockIncorrMessage[i]) != bvt.factBlockCorrectness[i] {
			isMatchedMes = "No"
			resultNotMatchedCounter++
		}
		log.Printf("   [%s %s %s %s", fmt.Sprint(i)+"]"+strings.Repeat(" ", 10-len(fmt.Sprint(i))), incorMes+strings.Repeat(" ", 44-len(incorMes)),
			verifMes+strings.Repeat(" ", 20-len(verifMes)), isMatchedMes)
	}

	result := "Passed"
	if bvt.blockAmmount-resultNotMatchedCounter != bvt.blockAmmount {
		result = "Not passed"
	}

	log.Printf("Test result: %d\\%d. %s\n", bvt.blockAmmount-resultNotMatchedCounter, bvt.blockAmmount, result)
}

func (bvt *BlockchainVerifTest) BlockchainVerefication() {
	bvt.source = rand.NewSource(time.Now().Unix())
	bvt.random = rand.New(bvt.source)
	blockchain.STARTBITS = 0xffff14
	bvt.genBlockTestAccounts(bvt.txPerBlockAmmount)
	log.Printf("Gnerated %d accounts", bvt.txPerBlockAmmount)
	log.Printf("Started generation of %d blocks (%d are incorrect)...", bvt.blockAmmount, bvt.incorrectblockAmmount)
	if bvt.incorrectblockAmmount > bvt.blockAmmount {
		log.Println("Wrong input")
	}
	bvt.genBlocks()
	bvt.printResults()
}
