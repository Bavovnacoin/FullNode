package blockchain

import (
	"bavovnacoin/byteArr"
	"bavovnacoin/hashing"
	"bavovnacoin/transaction"
	"bavovnacoin/txo"
	"bytes"
	"encoding/gob"
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

// Warning: it is considered that the block is valid
func AddBlockToBlockchain(block Block) bool {
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
	WriteBlock(BcLength, block)

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

func MineBlock(block Block, miningFlag int, allowPrint bool) Block {
	BlockForMining = block
	if miningFlag == 0 {
		block = MineThreads(block, 1, allowPrint)
	} else if miningFlag == 1 {
		block = MineThreads(block, uint64(runtime.NumCPU()), allowPrint)
	}

	return block
}

func VerifyBlock(block Block, height int, checkBits bool, allowCheckTxs bool) bool {
	var lastBlockHash string

	if int(BcLength) != 0 {
		var prevBlock Block

		if uint64(height) == BcLength {
			prevBlock = LastBlock
		} else {
			var isBlockFound bool
			prevBlock, isBlockFound = GetBlock(uint64(height) - 1)
			if !isBlockFound {
				return false
			}
		}

		if block.Time < prevBlock.Time {
			return false
		}

		lastBlockHash = hashing.SHA1(BlockHeaderToString(prevBlock))
	} else {
		lastBlockHash = "0000000000000000000000000000000000000000"
	}
	merkleRoot := GenMerkleRoot(block.Transactions)

	// Check bits value
	if checkBits {
		bits := GetCurrBitsValue()
		if bits != block.Bits {
			return false
		}
	}

	// Check nonce
	hashNonce, _ := new(big.Int).SetString(hashing.SHA1(BlockHeaderToString(block)), 16)
	if BitsToTarget(block.Bits).Cmp(hashNonce) != 1 {
		return false
	}

	// Check coinbase tx
	var allFee uint64
	for i := 1; i < len(block.Transactions); i++ {
		allFee += transaction.GetTxFee(block.Transactions[i])
	}
	if !CheckEmitedCoins(block.Transactions[0].Outputs[0].Value-allFee, height) {
		return false
	}

	// Check block hash values
	if block.HashPrevBlock != lastBlockHash ||
		block.MerkleRoot != merkleRoot {
		return false
	}

	// Check transactions
	if allowCheckTxs {
		for i := 1; i < int(len(block.Transactions)); i++ {
			if !transaction.VerifyTransaction(block.Transactions[i]) {
				return false
			}
		}
	}

	return true
}

func InitBlockchain() {
	BcLength, _ = GetBcHeight()
	if BcLength != 0 {
		LastBlock, _ = GetBlock(BcLength - 1)

		log.Println("Data is restored from db. Blockchain height:", BcLength)
	}
}

func FormGenesisBlock() Block {
	log.Println("Creating initial block")

	var rewardAdr byteArr.ByteArr
	rewardAdr.SetFromHexString(RewardAddress, 20)
	genesisBlock := CreateBlock(rewardAdr, true)
	genesisBlock.Bits = GetBits(true)
	genesisBlock = MineBlock(genesisBlock, 1, true)
	genesisBlock.Bits = STARTBITS

	if VerifyBlock(genesisBlock, int(BcLength), true, false) {
		AddBlockToBlockchain(genesisBlock)
		log.Println("Block is added to blockchain. Current height: " + fmt.Sprint(int(BcLength)+1) + "\n")

	} else {
		log.Println("Block is not added")
		println()
	}
	IncrBcHeight()
	return genesisBlock
}

func (block *Block) ToByteArr() ([]byte, bool) {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err := enc.Encode(block)
	if err != nil {
		log.Fatal("encode error:", err)
		return nil, false
	}

	return network.Bytes(), true
}
