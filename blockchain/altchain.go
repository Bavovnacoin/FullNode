package blockchain

import (
	"bavovnacoin/byteArr"
	"bavovnacoin/hashing"
	"bavovnacoin/transaction"
	"bavovnacoin/txo"
	"math/big"
)

var TransitionFactor = big.NewFloat(1.3)

func returnBlockTxo(block Block, height uint64) {
	for i := 0; i < len(block.Transactions); i++ {
		var txByteArr byteArr.ByteArr
		txByteArr.SetFromHexString(hashing.SHA1(transaction.GetCatTxFields(block.Transactions[i])), 20)

		outList := block.Transactions[i].Outputs
		for j := 0; j < len(outList); j++ {
			for k := 0; k < len(txo.CoinDatabase); k++ {
				if txo.CoinDatabase[k].OutTxHash.IsEqual(txByteArr) &&
					txo.CoinDatabase[k].TxOutInd == uint64(j) {
					txo.CoinDatabase = append(txo.CoinDatabase[:k], txo.CoinDatabase[k+1:]...)
					txo.RemUtxo(txByteArr, uint64(j), height)
					break
				}
			}
		}

		inpList := block.Transactions[i].Inputs
		for j := 0; j < len(inpList); j++ {
			txoToAdd, _ := txo.GetTxo(inpList[j].TxHash, inpList[j].OutInd, height)
			txo.RemoveTxo(inpList[j].TxHash, inpList[j].OutInd, height)

			txo.AddUtxo(txoToAdd.OutTxHash, txoToAdd.TxOutInd, txoToAdd.Value, txoToAdd.OutAddress, txoToAdd.BlockHeight)
			txo.CoinDatabase = append(txo.CoinDatabase, txoToAdd)
		}
	}
}

func addAltchBlockTxo(txs []transaction.Transaction) {
	for i := 0; i < len(txs); i++ {
		txInpList := txs[i].Inputs

		for j := 0; j < len(txInpList); j++ {
			txo.Spend(txInpList[j].TxHash, uint64(txInpList[j].OutInd))
		}

		txOutList := txs[i].Outputs
		for j := 0; j < len(txOutList); j++ {
			var txByteArr byteArr.ByteArr
			txByteArr.SetFromHexString(hashing.SHA1(transaction.GetCatTxFields(txs[i])), 20)
			txo.AddUtxo(txByteArr, uint64(j), txOutList[j].Value, txOutList[j].Address, BcLength-1)
		}
	}
}

/*
Reorganize procedure:
1. Store altchain and mainchain blocks; (DONE)
2. Rem utxo from mainchain (inputs -> (rem txo, add utxo), outputs -> remove utxo), add from altchain (DONE)
3. Delete them (do not delete from mainchain if altchain is higher then mainchain); (DONE)
4. Add under a new chainId and height; (DONE)
5. Swap heights in the end.

What do we do with the current mempool state (transactions in it are reffered to the old chain)
Nothing! - Solved by double spending check
*/
func reorganize(chainId uint64, altchHeight uint64) bool {
	blockIdToGet := altchHeight
	if altchHeight < BcLength {
		blockIdToGet = BcLength
	}

	forkHeight, _ := GetBlockForkHeight(chainId)
	for ; forkHeight <= blockIdToGet; forkHeight++ {
		// Get blocks to swap
		altChBlock, isAltGotten := GetBlock(forkHeight-1, chainId)
		mainChBlock, isMainGotten := GetBlock(forkHeight-1, 0)

		// Manage outputs and blocks
		if isMainGotten {
			returnBlockTxo(mainChBlock, forkHeight-1)
			RemBlock(forkHeight-1, 0)
		}

		if isAltGotten {
			addAltchBlockTxo(altChBlock.Transactions)
			RemBlock(forkHeight-1, chainId)
			println(len(txo.CoinDatabase))
		}

		if isMainGotten {
			WriteBlock(forkHeight-1, chainId, mainChBlock)
		}
		if isAltGotten {
			WriteBlock(forkHeight-1, 0, altChBlock)
			LastBlock = altChBlock
			BcLength++
		}
	}

	SetBcHeight(altchHeight, 0)
	SetBcHeight(BcLength, chainId)
	BcLength = altchHeight
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
	biggestChainIdHeight := uint64(0)
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
