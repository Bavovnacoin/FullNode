package blockchain

import (
	"bavovnacoin/hashing"
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
			block, _ := GetBlock(BcLength-i-1, 0)
			pastElevenBlocksSortedTime = util.InsertSorted(pastElevenBlocksSortedTime, block.Time)
		}

		pastElevenBlocksSortedTime = pastElevenBlocksSortedTime[:blocksToAddAmmount]
		timeCheckHeight = BcLength
	}
}

func checkCameBlockTime(blockTime int64, otherNodesTime []int64) bool {
	initPastBlocksTime()
	if blockTime < pastElevenBlocksSortedTime[len(pastElevenBlocksSortedTime)/2] ||
		(len(otherNodesTime) > 0 && blockTime > otherNodesTime[len(otherNodesTime)/2]+2) {
		return false
	}

	return true
}

func TryCameBlockToAdd(block Block, height uint64, otherNodesTime []int64) bool {
	PauseBlockAddition = true
	blockVer := !VerifyBlock(block, int(height), true, true)

	if blockVer || !checkCameBlockTime(block.Time, otherNodesTime) {
		PauseBlockAddition = false
		println("Came block is NOOTTT added!")
		return false
	}

	AllowMining = false
	BreakBlockAddition = true
	PauseBlockAddition = false
	var chainId uint64

	mainchBlock, _ := GetBlock(height-1, 0)

	var isAdded bool
	var blocks []BlockChainId
	if mainchBlock.Time >= mainchBlock.Time+1 && block.HashPrevBlock == hashing.SHA1(BlockHeaderToString(mainchBlock)) { // Fork from a mainchain
		blocks, _ = GetBlocksOnHeight(height)
		chainId = uint64(len(blocks))
		isAdded = WriteBlock(height, chainId, block)
		SetBlockForkHeight(height+1, chainId)
		SetBcHeight(height+1, chainId)
		println(1)
	} else { // Decide to what chain attach a new block
		blocks, chaindIds, _ := getAllLastBlocks()

		for i := 0; i < len(blocks); i++ {
			blockHash := hashing.SHA1(BlockHeaderToString(blocks[i]))

			if blockHash == block.HashPrevBlock {
				chainId = chaindIds[i]
				if chainId == 0 {
					isAdded = AddBlockToBlockchain(block, chainId, true)
					LastBlock = block
					IncrBcHeight(chainId)
					break
				} else {
					isAdded = WriteBlock(height, chainId, block)
					IncrBcHeight(chainId)
					break
				}
			}
		}
	}

	if TryReorganize() {
		log.Println("Reorganization happened")
	}

	if isAdded {
		if chainId == 0 {
			log.Println("Block is added to blockchain. Current height: " + fmt.Sprint(BcLength))
		}
		println("Came block is added!")
		println()
		return true
	}
	return false
}
