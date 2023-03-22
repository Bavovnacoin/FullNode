package synchronization

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/byteArr"
	"bavovnacoin/hashing"
)

type CheckPoint struct {
	height    uint64
	blockHash byteArr.ByteArr
}

var Checkpoints []CheckPoint
var checkpInd uint64

func createCheckpoint(height uint64, hashValue string) CheckPoint {
	var cp CheckPoint
	cp.height = height
	cp.blockHash.SetFromHexString(hashValue, 20)
	return cp
}

func InitCheckpoints() {
	checkpInd = 0
	Checkpoints = append(Checkpoints, createCheckpoint(10, ""))
}

func checkForCheckpCorrespondence(height uint64, block blockchain.Block) bool {
	if height == Checkpoints[checkpInd].height {
		var blockHash byteArr.ByteArr
		blockHash.SetFromHexString(hashing.SHA1(blockchain.BlockToString(block)), 20)

		if blockHash.IsEqual(Checkpoints[checkpInd].blockHash) {
			checkpInd++
			return true
		}
	}
	return false
}
