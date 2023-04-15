package blockchain

import (
	"bavovnacoin/byteArr"
	"bavovnacoin/hashing"
	"bavovnacoin/transaction"
	"bavovnacoin/txo"
	"fmt"
	"log"
	"math/big"
	"runtime"
	"strings"
	"time"
)

var BcLength uint64
var LastBlock Block
var IsMempAdded bool
var BlockForMining Block

var CreatedBlock Block
var AllowCreateBlock bool = true

var PauseBlockAddition bool
var BreakBlockAddition bool

var RewardAddress string = "9e90c94ab3b2da7900bdc70680f4a9c8f2fe0375"

func getChainwork(block Block) *big.Int {
	return new(big.Int).Add(LastBlock.Chainwork, getCurrBlockChainwork(block))
}

// Warning: it is considered that the block is valid
func AddBlockToBlockchain(block Block, chainId uint64) bool {
	for PauseBlockAddition {
		time.Sleep(10 * time.Millisecond)
		if BreakBlockAddition {
			BreakBlockAddition = false
			return false
		}
	}

	for i := 0; i < len(block.Transactions); i++ {
		txInpList := block.Transactions[i].Inputs

		for j := 0; j < len(txInpList); j++ {
			txo.Spend(txInpList[j].TxHash, uint64(txInpList[j].OutInd))
		}

		txOutList := block.Transactions[i].Outputs
		for j := 0; j < len(txOutList); j++ {
			var txByteArr byteArr.ByteArr
			txByteArr.SetFromHexString(hashing.SHA1(transaction.GetCatTxFields(block.Transactions[i])), 20)
			txo.AddUtxo(txByteArr, uint64(j), txOutList[j].Value, txOutList[j].Address, uint64(int(BcLength)))
		}
	}

	LastBlock = block
	WriteBlock(BcLength, chainId, block)

	IsMempAdded = false
	return true
}

func GetBits(allowPrint bool) uint64 {
	bits := GetCurrBitsValue()

	if allowPrint {
		target := fmt.Sprintf("%x", BitsToTarget(bits))
		log.Println("Current bits value is " + fmt.Sprintf("%x", bits))
		log.Println("Current target value is " + strings.Repeat("0", 40-len(target)) + target)
	}
	return bits
}

func MineBlock(block Block, miningFlag int, allowPrint bool) (Block, bool) {
	BlockForMining = block
	miningRes := true
	if miningFlag == 0 {
		block, miningRes = MineThreads(block, 1, allowPrint)
	} else if miningFlag == 1 {
		block, miningRes = MineThreads(block, uint64(runtime.NumCPU()), allowPrint)
	}

	return block, miningRes
}

func VerifyBlock(block Block, height int, checkBits bool, allowCheckTxs bool) bool {
	var lastBlockHashes []string
	if int(BcLength) != 0 {
		var prevBlocks []Block

		if uint64(height) == BcLength {
			prevBlocks = append(prevBlocks, LastBlock)
		} else {
			var isBlockFound bool
			prevBlocks, isBlockFound = GetBlocksOnHeight(uint64(height) - 1)
			if !isBlockFound {
				println("Block found problem")
				return false
			}
		}

		for i := 0; i < len(prevBlocks); i++ {
			lastBlockHashes = append(lastBlockHashes, hashing.SHA1(BlockHeaderToString(prevBlocks[i])))
		}
	} else {
		lastBlockHashes = append(lastBlockHashes, "0000000000000000000000000000000000000000")
	}
	merkleRoot := GenMerkleRoot(block.Transactions)

	// Check bits value
	if checkBits {
		bits := GetCurrBitsValue()
		if bits != block.Bits {
			println("Bits problem")
			return false
		}
	}

	// Check nonce
	h := hashing.SHA1(BlockHeaderToString(block))
	hashNonce, _ := new(big.Int).SetString(h, 16)
	if BitsToTarget(block.Bits).Cmp(hashNonce) == -1 {
		fmt.Println(hashNonce.Bytes())
		println(fmt.Sprintf("%x", hashNonce), fmt.Sprintf("%x", BitsToTarget(block.Bits)))
		println("Nonce problem", hashing.SHA1(BlockHeaderToString(block)), fmt.Sprintf("%x", BitsToTarget(block.Bits)))
		return false
	}

	// Check coinbase tx
	var allFee uint64
	for i := 1; i < len(block.Transactions); i++ {
		allFee += transaction.GetTxFee(block.Transactions[i])
	}
	if !CheckEmitedCoins(block.Transactions[0].Outputs[0].Value-allFee, height) {
		println("Emited coins problem")
		return false
	}

	// Check block hash values
	var hashFound bool
	for i := 0; i < len(lastBlockHashes); i++ {
		if block.HashPrevBlock == lastBlockHashes[i] {
			hashFound = true
			break
		}
	}
	if !hashFound {
		println("Hash problem")
		return false
	}

	// Check Merkle root
	if block.MerkleRoot != merkleRoot {
		println("Merkle root problem")
		return false
	}

	// Check transactions
	if allowCheckTxs {
		for i := 1; i < int(len(block.Transactions)); i++ {
			if !transaction.VerifyTransaction(block.Transactions[i]) {
				println("Tx problem")
				return false
			}
		}
	}

	// Check chainwork
	chainwork := getChainwork(block)
	if block.Chainwork.Cmp(chainwork) != 0 {
		return false
	}

	return true
}

func InitBlockchain() {
	BcLength, _ = GetBcHeight()
	if BcLength != 0 {
		LastBlock, _ = GetBlock(BcLength-1, 0)

		log.Println("Data is restored from db. Blockchain height:", BcLength)
	}
}

func FormGenesisBlock() Block {
	log.Println("Creating initial block")

	var rewardAdr byteArr.ByteArr
	rewardAdr.SetFromHexString(RewardAddress, 20)
	genesisBlock := CreateBlock(rewardAdr, true)
	genesisBlock.Bits = GetBits(true)
	genesisBlock, _ = MineBlock(genesisBlock, 1, true)
	genesisBlock.Bits = STARTBITS

	if VerifyBlock(genesisBlock, int(BcLength), true, false) {
		AddBlockToBlockchain(genesisBlock, 0)
		log.Println("Block is added to blockchain. Current height: " + fmt.Sprint(int(BcLength)+1) + "\n")
	} else {
		log.Println("Block is not added")
		println()
	}
	IncrBcHeight()
	return genesisBlock
}

func IsBlockExists(blockHash byteArr.ByteArr, height uint64) bool {
	println("Checking block hash", blockHash.ToHexString(), height)

	blockArr, res := GetBlocksOnHeight(height)
	if !res {
		println("Block with such a height is not found")
		return false
	}

	var bcBlockHash byteArr.ByteArr
	for i := 0; i < len(blockArr); i++ {
		bcBlockHash.SetFromHexString(hashing.SHA1(BlockHeaderToString(blockArr[i])), 20)
		if bcBlockHash.IsEqual(blockHash) {
			println("Block is found")
			return true
		}
	}

	println("Block is not found", bcBlockHash.ToHexString(), blockHash.ToHexString())
	return false
}
