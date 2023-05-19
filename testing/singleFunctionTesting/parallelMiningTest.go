/*
	Checks how is mining performed with help of parallel computations
*/

package singleFunctionTesting

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/hashing"
	"bavovnacoin/node/node_settings"
	"fmt"
	"log"
	"math/big"
	"runtime"
)

type ParallelMiningTest struct {
	ThreadsCount int
	TestsAmmount int
	Bits         int
}

func (pmt *ParallelMiningTest) mineBlocks() {
	block := blockchain.Block{}
	var miningRes bool
	block.Bits = 0xffff13
	var passedTestAmmount int

	log.Println("Test started")
	for i := 0; i < pmt.TestsAmmount; i++ {
		block.HashPrevBlock = hashing.SHA1(blockchain.BlockHeaderToString(block))
		block, miningRes = blockchain.MineThreads(block, false)

		blockHash, _ := new(big.Int).SetString(hashing.SHA1(blockchain.BlockHeaderToString(block)), 16)
		if miningRes && blockchain.BitsToTarget(block.Bits).Cmp(blockHash) == 1 {
			fmt.Printf("[%d]. Passed. Nonce value: %d\n", i+1, block.Nonce)
			passedTestAmmount++
		} else {
			fmt.Printf("[%d]. Incorrect mining. Test string: %s. Bits value: %x\n", i+1, blockchain.BlockHeaderToString(block), block.Bits)
		}
	}

	if passedTestAmmount == pmt.TestsAmmount {
		log.Printf("Test passed (%d/%d)!\n", passedTestAmmount, pmt.TestsAmmount)
	} else {
		log.Printf("Test is not passed (%d/%d)!\n", passedTestAmmount, pmt.TestsAmmount)
	}
}

func (pmt *ParallelMiningTest) Launch(testAmmount int) {
	if pmt.ThreadsCount <= 1 {
		pmt.ThreadsCount = runtime.NumCPU()
	}

	pmt.TestsAmmount = testAmmount

	node_settings.Settings.MiningThreads = uint(pmt.ThreadsCount)
	pmt.mineBlocks()
}
