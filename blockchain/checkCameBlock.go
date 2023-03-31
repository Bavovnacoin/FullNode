package blockchain

import (
	"bavovnacoin/util"
	"fmt"
	"log"
)

var pastElevenBlocksSortedTime []int64
var timeCheckHeight uint64

func initPastBlocksTime() {
	if timeCheckHeight != BcLength {
		blocksToAddAmmount := BcLength - timeCheckHeight
		if blocksToAddAmmount > 11 {
			blocksToAddAmmount = 11
		}

		for i := uint64(0); i < blocksToAddAmmount; i++ {
			block, _ := GetBlock(BcLength - i - 1)
			pastElevenBlocksSortedTime = util.InsertSorted(pastElevenBlocksSortedTime, block.Time)
		}

		pastElevenBlocksSortedTime = pastElevenBlocksSortedTime[:blocksToAddAmmount]
		timeCheckHeight = BcLength
	}
}

func checkCameBlockTime(blockTime int64, otherNodesTime []int64) bool {
	initPastBlocksTime()
	if blockTime < pastElevenBlocksSortedTime[len(pastElevenBlocksSortedTime)/2] ||
		blockTime > otherNodesTime[len(otherNodesTime)/2]+2 {
		return false
	}

	return true
}

func TryCameBlockToAdd(block Block, otherNodesTime []int64) bool {
	println("Trying to add a new block")
	PauseBlockAddition = true
	blockVer := !VerifyBlock(block, int(BcLength), true, true)
	println("Block verefication completed")
	if blockVer || !checkCameBlockTime(block.Time, otherNodesTime) {
		PauseBlockAddition = false
		println("Came block is NOOTTT added!")
		return false
	}
	AllowMining = false
	BreakBlockAddition = true
	PauseBlockAddition = false
	println("New block checked")

	// TODO: solve disagreement (equal time change to block created eaarlier)
	if CreatedBlock.Time == block.Time {

	}

	AddBlockToBlockchain(block)
	BreakBlockAddition = false
	println("New block Added")
	log.Println("Block is added to blockchain. Current height: " + fmt.Sprint(BcLength+1) + "\n")
	println()
	IncrBcHeight()

	AllowCreateBlock = true
	CreatedBlock.MerkleRoot = ""
	println("Came block is added!")
	return true
}
