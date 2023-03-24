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
	Checkpoints = append(Checkpoints, setCheckpoint(2, "0529fe2df59de3ed8a9b1c3c75746d735a161afc"))
	Checkpoints = append(Checkpoints, setCheckpoint(5, "6ce9d14e492b66c785c8ac4199aade93e7403049"))
}

func checkForCheckpCorrespondence(height uint64, block blockchain.Block) bool {
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
