package singleFunctionTesting

import (
	"bavovnacoin/account"
	"bavovnacoin/blockchain"
	"bavovnacoin/byteArr"
	"bavovnacoin/dbController"
	"bavovnacoin/hashing"
	"bavovnacoin/node"
	"bavovnacoin/node_controller/command_executor"
	"bavovnacoin/transaction"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type ReorganizationVerifTest struct {
	SingleFunctionTesting
	mcBlockAmmount uint64
	acBlockAmmount uint64
	prevHeight     uint64
	currAccInd     int

	source rand.Source
	random *rand.Rand
}

func (rv *ReorganizationVerifTest) CreateTx() transaction.Transaction {
	fee := rv.random.Intn(5) + 1
	isGenLocktime := rv.random.Intn(5)
	var locktime uint
	if isGenLocktime == 2 {
		locktime = uint(int(blockchain.BcLength+1) + rv.random.Intn(2) + 1)
	}

	var outAddr []byteArr.ByteArr
	outAddr = append(outAddr, byteArr.ByteArr{})
	outAddr[0].SetFromHexString(hashing.SHA1(account.Wallet[rv.currAccInd].KeyPairList[0].PublKey), 20)

	var outValue []uint64
	outValue = append(outValue, account.CurrAccount.Balance/2)

	tx, _ := transaction.CreateTransaction(fmt.Sprint(rv.currAccInd), outAddr, outValue, fee, locktime)
	rv.currAccInd++
	return tx
}

func (rv *ReorganizationVerifTest) nodeWorkListener(blocksCount uint64) {
	for true {
		if blockchain.BcLength >= blocksCount {
			command_executor.ComContr.FullNodeWorking = false
			return
		}

		if rv.prevHeight != blockchain.BcLength {
			rv.prevHeight = blockchain.BcLength
			blockchain.AddTxToMempool(rv.CreateTx(), false)
		}
	}
}

func (rv *ReorganizationVerifTest) genBlocks() {
	command_executor.ComContr.FullNodeWorking = true
	go rv.nodeWorkListener(rv.mcBlockAmmount)
	node.BlockGen(false)
}

func (rv *ReorganizationVerifTest) genAltchBlocks() {
	bl, _ := blockchain.GetBlock(blockchain.BcLength-2, 0)
	var prevHash string = hashing.SHA1(blockchain.BlockHeaderToString(bl))

	for i := uint64(0); i < rv.acBlockAmmount; i++ {
		blockchain.AddTxToMempool(rv.CreateTx(), false)
		node.CreateBlockLog(blockchain.GetBits(true), prevHash, bl, true)
		blockchain.AllowCreateBlock = false

		var otherNodesTime []int64
		otherNodesTime = append(otherNodesTime, time.Now().UTC().Unix())
		if i == 0 {
			blockchain.CreatedBlock.Time = bl.Time
		}
		blockchain.CreatedBlock.Version = 1
		blockchain.TryCameBlockToAdd(blockchain.CreatedBlock, blockchain.BcLength-1+uint64(i), otherNodesTime)
		prevHash = hashing.SHA1(blockchain.BlockHeaderToString(blockchain.CreatedBlock))
		bl = blockchain.CreatedBlock
	}
}

func (rv *ReorganizationVerifTest) printResult() {
	println("Results:")
	println("Blockchain scheme:")
	for height := 0; true; height++ {
		blocks, res := blockchain.GetBlocksOnHeight(uint64(height))
		if !res || len(blocks) == 0 {
			break
		}

		var str string
		if len(blocks) == 1 {
			if blocks[0].ChainId == 0 {
				str += fmt.Sprint(blocks[0].Block.Version)
				str += "  "
			} else {
				str += "  "
				str += fmt.Sprint(blocks[0].Block.Version)
			}
		} else {
			str += fmt.Sprintf("%s %s", fmt.Sprint(blocks[0].Block.Version), fmt.Sprint(blocks[1].Block.Version))
		}

		println(str)
	}

}

func (rv *ReorganizationVerifTest) Launch() {
	rv.mcBlockAmmount = 2
	rv.acBlockAmmount = 3

	dbController.DbPath = "testing/testData"
	if _, err := os.Stat(dbController.DbPath); err == nil {
		os.RemoveAll(dbController.DbPath)
		println("Removed test db from a previous test.")
	}
	dbController.DB.OpenDb()

	blockchain.STARTBITS = 0xffff14
	rv.source = rand.NewSource(time.Now().Unix())
	rv.random = rand.New(rv.source)

	rv.genBlockTestAccounts(int(rv.mcBlockAmmount) + int(rv.acBlockAmmount))
	rv.genBlocks() // Generating blocks in mainchain

	rv.genAltchBlocks()

	rv.printResult()
}
