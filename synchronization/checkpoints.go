package synchronization

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/byteArr"
	"bavovnacoin/dbController"
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
	// Checkpoints = append(Checkpoints, setCheckpoint(2, "00051ebaef762eea0932bb97d61c76061899868d"))
	// Checkpoints = append(Checkpoints, setCheckpoint(5, "000ab237e87118a9c4077ae20854b304d0ca8424"))
}

func checkForBlockCorrespondence(height uint64, block blockchain.Block) bool {
	if checkpInd < uint64(len(Checkpoints)) && height == Checkpoints[checkpInd].height {
		var blockHash byteArr.ByteArr
		blockHash.SetFromHexString(hashing.SHA1(blockchain.BlockHeaderToString(block)), 20)
		if blockHash.IsEqual(Checkpoints[checkpInd].blockHash) {
			checkpInd++
			return true
		} else {
			return false
		}
	} else if checkpInd >= uint64(len(Checkpoints)) {
		isVal := blockchain.VerifyBlock(block, int(height), true, true)
		return isVal
	}

	return true
}

func GetCheckpHashes(args ...uint64) {
	dbController.DB.OpenDb()
	defer dbController.DB.CloseDb()
	blockchain.InitBlockchain()

	for _, ind := range args {
		b, _ := blockchain.GetBlock(ind, 0)
		println(hashing.SHA1(blockchain.BlockHeaderToString(b)))
	}
}
