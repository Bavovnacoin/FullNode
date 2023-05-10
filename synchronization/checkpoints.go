package synchronization

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/byteArr"
	"bavovnacoin/dbController"
	"bavovnacoin/hashing"
)

type CheckPoint struct {
	Height    uint64
	BlockHash byteArr.ByteArr
}

var Checkpoints []CheckPoint
var checkpInd uint64

func SetCheckpoint(height uint64, hashValue string) CheckPoint {
	var cp CheckPoint
	cp.Height = height
	cp.BlockHash.SetFromHexString(hashValue, 20)
	return cp
}

func InitCheckpoints() {
	checkpInd = 0
	//Checkpoints = append(Checkpoints, setCheckpoint(1, "00004b3e6648b1ccaa64f45634d0c2d15c7e7c02"))
}

func CheckForBlockCorrespondence(height uint64, block blockchain.Block) bool {
	if blockchain.VerifyBlock(block, int(height), true, true) {
		if checkpInd < uint64(len(Checkpoints)) && height == Checkpoints[checkpInd].Height {
			var blockHash byteArr.ByteArr
			blockHash.SetFromHexString(hashing.SHA1(blockchain.BlockHeaderToString(block)), 20)
			if blockHash.IsEqual(Checkpoints[checkpInd].BlockHash) {
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
