package blockchain

import (
	"bavovnacoin/util"
)

var pastElevenBlocksSortedTime []int64
var timeCheckHeight uint64

func initPastBlocksTime() {
	if timeCheckHeight != BcLength {
		blocksToAddAmmount := BcLength - timeCheckHeight
		if blocksToAddAmmount > 11 {
			blocksToAddAmmount = 11
		}

		for i := BcLength - 1; i >= 0 || i >= BcLength-blocksToAddAmmount; i-- {
			block, _ := GetBlock(i)
			pastElevenBlocksSortedTime = util.InsertSorted(pastElevenBlocksSortedTime, block.Time)
		}

		pastElevenBlocksSortedTime = pastElevenBlocksSortedTime[:11]
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

func AllowCameBlockToAdd(block Block, otherNodesTime []int64) bool {
	if !checkCameBlockTime(block.Time, otherNodesTime) || !VerifyBlock(block, int(BcLength)+1, true, true) {
		return false
	}
	return true
}
