package singleFunctionTesting

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/dbController"
	"bavovnacoin/hashing"
	"bavovnacoin/node"
	"bavovnacoin/node_controller/command_executor"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type ReorganizationVerifTest struct {
	SingleFunctionTesting
	blockAmmount uint64

	source rand.Source
	random *rand.Rand
}

func (rv *ReorganizationVerifTest) nodeWorkListener(blocksCount uint64) {
	for true {
		if blockchain.BcLength >= blocksCount {
			command_executor.ComContr.FullNodeWorking = false
			return
		}
	}
}

func (rv *ReorganizationVerifTest) genBlocks() {
	command_executor.ComContr.FullNodeWorking = true
	go rv.nodeWorkListener(rv.blockAmmount)
	node.BlockGen(false)
}

func (rv *ReorganizationVerifTest) genAltchBlocks() {
	bl, _ := blockchain.GetBlock(blockchain.BcLength-2, 0)
	blockchain.PrintBlockTitle(bl, blockchain.BcLength-2)
	var prevHash string = hashing.SHA1(blockchain.BlockHeaderToString(bl))

	for i := uint64(0); i < rv.blockAmmount; i++ {
		node.CreateBlockLog(blockchain.GetBits(true), prevHash, bl, true)
		blockchain.AllowCreateBlock = false

		var otherNodesTime []int64
		otherNodesTime = append(otherNodesTime, time.Now().UTC().Unix())
		if i == 0 {
			blockchain.CreatedBlock.Time = bl.Time
			//blockchain.CreatedBlock.Chainwork = bl.Chainwork
		}
		blockchain.CreatedBlock.Version = 1
		blockchain.TryCameBlockToAdd(blockchain.CreatedBlock, blockchain.BcLength-1+uint64(i), otherNodesTime)
		prevHash = hashing.SHA1(blockchain.BlockHeaderToString(blockchain.CreatedBlock))
		bl = blockchain.CreatedBlock
	}
}

func (rv *ReorganizationVerifTest) joinChainToString(arr []string) string {
	var res string

	for i := 0; i < len(arr); i++ {
		res += arr[i]
		if arr[i] != " " {
			if i < len(arr)-1 {
				res += "-"
			}
		} else {
			res += " "

		}
	}
	return res
}

func (rv *ReorganizationVerifTest) printResult() {
	//TODO: check TXO
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
				//fmt.Printf("%d\n", blocks[0].Block.Chainwork)
			}
		} else {
			str += fmt.Sprintf("%s %s", fmt.Sprint(blocks[0].Block.Version), fmt.Sprint(blocks[1].Block.Version))
			//fmt.Printf("%d\n", blocks[1].Block.Chainwork)
		}

		println(str)
	}

}

func (rv *ReorganizationVerifTest) Launch() {
	rv.blockAmmount = 3

	dbController.DbPath = "testing/testData"
	if _, err := os.Stat(dbController.DbPath); err == nil {
		os.RemoveAll(dbController.DbPath)
		println("Removed test db from a previous test.")
	}
	dbController.DB.OpenDb()

	blockchain.STARTBITS = 0xffff14
	rv.source = rand.NewSource(time.Now().Unix())
	rv.random = rand.New(rv.source)

	rv.genBlockTestAccounts(int(rv.blockAmmount))
	rv.genBlocks() // Generating blocks in mainchain

	rv.genAltchBlocks()

	// TODO: gen last block, so it can be verified that it added to the new mainchain
	rv.printResult()
}
