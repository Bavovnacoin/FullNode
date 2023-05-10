/*
	Tests how blocks are downloaded and uploaded (including altchains)
*/

package singleFunctionTesting

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/dbController"
	"bavovnacoin/hashing"
	"bavovnacoin/networking_p2p"
	"bavovnacoin/synchronization"
	"os"
)

type SynchronizationTest struct {
	SingleFunctionTesting

	mcBlocksAmmount int
	acBlocksAmmount int
	acAmmount       int
}

func (st SynchronizationTest) genAltchain(chainId uint64) {
	startHeight := uint64(st.mcBlocksAmmount) / 2
	prevBlock, _ := blockchain.GetBlock(startHeight, 0)

	newBlock := blockchain.Block{}
	for i := 0; i < st.acBlocksAmmount; i++ {
		newBlock, _ = st.CreateBlock(blockchain.GetBits(false),
			hashing.SHA1(blockchain.BlockHeaderToString(prevBlock)), prevBlock, false)
		newBlock.Version = uint(chainId)

		blockchain.TryCameBlockToAdd(newBlock, startHeight+uint64(i)+1, []int64{}, false)
		prevBlock = newBlock
	}
}

func (st SynchronizationTest) genMainchain() {
	newBlock := blockchain.Block{}
	for i := 0; i < st.mcBlocksAmmount; i++ {
		newBlock, _ = st.CreateBlock(blockchain.GetBits(false),
			hashing.SHA1(blockchain.BlockHeaderToString(blockchain.LastBlock)), blockchain.LastBlock, false)

		if i == 0 {
			newBlock.HashPrevBlock = "0000000000000000000000000000000000000000"
		}

		blockchain.AddBlockToBlockchain(newBlock, blockchain.BcLength, 0, true)
		blockchain.IncrBcHeight(0)
		blockchain.LastBlock = newBlock

		if i%5 == 0 {
			synchronization.Checkpoints = append(synchronization.Checkpoints, synchronization.SetCheckpoint(uint64(i), hashing.SHA1(blockchain.BlockHeaderToString(newBlock))))
		}
	}
}

func (st SynchronizationTest) genChains() {
	st.genMainchain()
	for i := 0; i < st.acAmmount; i++ {
		st.genAltchain(uint64(i + 1))
	}
}

func (st SynchronizationTest) startSync() (bool, [][]networking_p2p.BlocksOnHeight) {
	var blocksToAdd [][]networking_p2p.BlocksOnHeight
	var blockReq networking_p2p.BlockRequest
	blockReq.IsMoreBlocks = true
	blockReq.Blocks = append(blockReq.Blocks, networking_p2p.BlocksOnHeight{Height: 0})
	reqHeight := uint64(0)

	for blockReq.IsMoreBlocks {
		blockReq = networking_p2p.GetBlocksOnHeight(reqHeight)
		for j := 0; j < len(blockReq.Blocks); j++ {
			reqHeight += uint64(len(blockReq.Blocks[j].Blocks))
		}

		blocksToAdd = append(blocksToAdd, blockReq.Blocks)
	}

	dbController.DB.CloseDb()

	os.RemoveAll(dbController.DbPath)
	blockchain.BcLength = 0
	InitTestDb(false)

	var syncRes bool
	for i := 0; i < len(blocksToAdd); i++ {
		syncRes = networking_p2p.SyncAddBlocks(blocksToAdd[i], true)
		if !syncRes {
			break
		}
	}
	return syncRes, blocksToAdd
}

func (st SynchronizationTest) printRes(syncRes bool, genChain [][]networking_p2p.BlocksOnHeight) {
	var checkRes bool = true

	for i := uint64(0); i < uint64(len(genChain)); i++ {
		for j := 0; j < len(genChain[i]); j++ {
			blocks, _ := blockchain.GetBlocksOnHeight(genChain[i][j].Height)
			for k := 0; k < len(blocks); k++ {
				if hashing.SHA1(blockchain.BlockHeaderToString(genChain[i][j].Blocks[k].Block)) !=
					hashing.SHA1(blockchain.BlockHeaderToString(blocks[k].Block)) {
					checkRes = false
					break
				}
			}
		}
	}

	if checkRes == syncRes && syncRes == true {
		println("Test passed!")
	} else {
		println("Test is not passed!")
	}
}

func (st SynchronizationTest) Launch() {
	st.mcBlocksAmmount = 10
	st.acBlocksAmmount = 3
	st.acAmmount = 1
	blockchain.STARTBITS = 0xffff14

	InitTestDb(true)
	st.genChains()
	syncRes, genChain := st.startSync()

	st.printRes(syncRes, genChain)
}
