/*
	Checks whether difficulty is changed and whether it changed correctly
*/

package singleFunctionTesting

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/hashing"
	"fmt"
	"time"
)

type DifficultyChangeTest struct {
	SingleFunctionTesting
	StartBits uint64
}

func (dcht *DifficultyChangeTest) genBlocks(isLong bool) bool {
	blockchain.LastBlock.Time = time.Now().Unix()
	for i := 0; i < blockchain.BLOCK_DIFF_CHECK; i++ {
		newBlock, _ := dcht.CreateBlock(blockchain.GetBits(false),
			hashing.SHA1(blockchain.BlockHeaderToString(blockchain.LastBlock)), blockchain.LastBlock, false)

		if isLong {
			newBlock.Time = blockchain.LastBlock.Time + int64(blockchain.BLOCK_CREATION_SEC) + 2
		}

		blockchain.AddBlockToBlockchain(blockchain.CreatedBlock, blockchain.BcLength, 0, true)
		blockchain.IncrBcHeight(0)
		blockchain.LastBlock = newBlock
	}

	genBits := blockchain.GetBits(false)
	isIncr := true
	fmt.Printf("Changed bits value from %x to %x ", dcht.StartBits, genBits)

	if blockchain.BitsToTarget(dcht.StartBits).Cmp(blockchain.BitsToTarget(genBits)) == -1 {
		fmt.Printf("(increased by %d times)\n", (genBits / dcht.StartBits))
	} else if blockchain.BitsToTarget(dcht.StartBits).Cmp(blockchain.BitsToTarget(genBits)) == 1 {
		fmt.Printf("(reduced by %d times)\n", (genBits / dcht.StartBits))
		isIncr = false
	}

	if (isLong && isIncr) || (!isLong && !isIncr) {
		return true
	}

	return false
}

func (dcht *DifficultyChangeTest) Launch() {
	blockchain.BLOCK_DIFF_CHECK = 10
	dcht.StartBits = 0x00ffff13
	blockchain.STARTBITS = dcht.StartBits

	InitTestDb(true)

	println("Generating bc faster than expected")
	littleTimeRes := dcht.genBlocks(false)
	blockchain.LastBlock = blockchain.Block{}
	blockchain.BcLength = 0
	println()

	InitTestDb(false)
	println("Generating bc slower than expected")
	longTimeRes := dcht.genBlocks(true)

	if littleTimeRes && longTimeRes {
		println("Test passed!")
	} else {
		println("Test is not passed!")
	}
}
