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

func setCheckpoint(height uint64, hashValue string) CheckPoint {
	var cp CheckPoint
	cp.height = height
	cp.blockHash.SetFromHexString(hashValue, 20)
	return cp
}

func InitCheckpoints() {
	checkpInd = 0
	Checkpoints = append(Checkpoints, setCheckpoint(2, "0013ef434f366f596593738a9e50ddd5ea06ca66"))
	Checkpoints = append(Checkpoints, setCheckpoint(5, "00127ebf6d909a17be9e41422589b06d1f607de9"))
}

func checkForBlockCorrespondence(height uint64, block blockchain.Block) bool {
	if checkpInd < uint64(len(Checkpoints)) && height == Checkpoints[checkpInd].height {
		var blockHash byteArr.ByteArr
		blockHash.SetFromHexString(hashing.SHA1(blockchain.BlockToString(block)), 20)

		println(blockHash.ToHexString(), Checkpoints[checkpInd].blockHash.ToHexString())
		if blockHash.IsEqual(Checkpoints[checkpInd].blockHash) {
			checkpInd++
			return true
		} else {
			return false
		}
	} else if checkpInd >= uint64(len(Checkpoints)) {
		isVal := blockchain.VerifyBlock(block, int(height), true, true)
		println(isVal, height)
		return isVal
	}

	return true
}
