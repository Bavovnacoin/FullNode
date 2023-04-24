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
	Checkpoints = append(Checkpoints, setCheckpoint(1, "00004b3e6648b1ccaa64f45634d0c2d15c7e7c02"))
}

func checkForBlockCorrespondence(height uint64, block blockchain.Block) bool {
	if blockchain.VerifyBlock(block, int(height), true, true) {
		if checkpInd < uint64(len(Checkpoints)) && height == Checkpoints[checkpInd].height {
			var blockHash byteArr.ByteArr
			blockHash.SetFromHexString(hashing.SHA1(blockchain.BlockHeaderToString(block)), 20)
			if blockHash.IsEqual(Checkpoints[checkpInd].blockHash) {
				checkpInd++
				return true
			} else {
				return false
			}
		}
		return true
	}

	return false
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
