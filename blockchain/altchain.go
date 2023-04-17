package blockchain

import "math/big"

var TransitionFactor = big.NewFloat(1.5)

/*
	Reorganize procedure:
	1. store altchain and mainchain blocks; (DONE)
	2. Rem utxo from mainchain, add from altchain
	3. delete them (do not delete from mainchain if altchain is higher then mainchain);
	4. add under a new chainId and height;
	5. swap heights in the end.

	What do we do with the current mempool state (transactions in it are reffered to the old chain)
	Nothing!
*/
func reorganize(chainId uint64, height uint64) bool {
	// Use "forkheight"
	// for true {
	// 	altchBlock, isAltBlockGotten := GetBlock(height-1, chainId)
	// 	if !isAltBlockGotten {
	// 		return true
	// 	}

	// 	mainchBlock, isMnBlockGotten := GetBlock(height-1, 0)

	// 	height--
	// }
	return true
}

func TryReorganize() bool {
	lastBlocks, chainIds, heights := getAllLastBlocks()
	var mainchainArrId int = -1

	// Select mainchain from all chains last blocks
	for i := 0; i < len(chainIds); i++ {
		if chainIds[i] == 0 {
			mainchainArrId = i
			break
		}
	}

	if mainchainArrId == -1 {
		return false
	}

	biggestFact := TransitionFactor
	biggestFactChainId := uint64(0)
	var biggestChainIdHeight uint64
	allowReorg := false

	mainchWork := new(big.Float).SetInt(lastBlocks[mainchainArrId].Chainwork)

	// Compare all chains to a transition factor
	for i := 0; i < len(lastBlocks); i++ {
		if chainIds[i] == 0 {
			continue
		}

		fact := new(big.Float).Quo(new(big.Float).SetInt(lastBlocks[i].Chainwork), mainchWork)
		if fact.Cmp(biggestFact) == 1 {
			biggestFact = fact
			biggestFactChainId = chainIds[i]
			biggestChainIdHeight = heights[i]
			allowReorg = true
		}
	}

	// Reorganize if it is a need
	var reorgRes bool
	if allowReorg {
		reorgRes = reorganize(biggestFactChainId, biggestChainIdHeight)
	}

	return reorgRes
}
