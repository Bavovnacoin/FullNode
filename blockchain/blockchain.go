package blockchain

import (
	"bavovnacoin/byteArr"
	"bavovnacoin/hashing"
	"bavovnacoin/transaction"
	"bavovnacoin/utxo"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"math/big"
	"runtime"
	"strings"
	"time"
)

var BcLength uint64
var LastBlock Block
var IsMempAdded bool
var BlockForMining Block

type Block struct {
	Blocksize        uint
	Version          uint
	HashPrevBlock    string
	Time             time.Time
	TransactionCount uint
	MerkleRoot       string
	Bits             uint64
	Nonce            uint64
	Transactions     []transaction.Transaction
}

func BlockToString(block Block) string {
	str := ""
	str += fmt.Sprint(block.Blocksize)
	str += fmt.Sprint(block.Version)
	str += block.HashPrevBlock
	str += block.Time.String()
	str += fmt.Sprint(block.TransactionCount)
	str += block.MerkleRoot
	str += fmt.Sprintf("%x", block.Bits)
	str += fmt.Sprint(block.Nonce)
	for i := 0; i < len(block.Transactions); i++ {
		str += transaction.GetCatTxFields(block.Transactions[i])
	}
	return str
}

func AddBlockToBlockchain(block Block) bool {
	isBlockValid := ValidateBlock(block, int(BcLength))

	if isBlockValid {
		for i := 0; i < len(block.Transactions); i++ {
			txInpList := block.Transactions[i].Inputs

			for j := 0; j < len(txInpList); j++ {
				utxo.Spend(txInpList[j].TxHash, uint64(txInpList[j].OutInd))
			}

			txOutList := block.Transactions[i].Outputs
			for j := 0; j < len(txOutList); j++ {
				var txByteArr byteArr.ByteArr
				txByteArr.SetFromHexString(hashing.SHA1(transaction.GetCatTxFields(block.Transactions[i])), 20)
				utxo.AddUtxo(txByteArr, uint64(j), txOutList[j].Value, txOutList[j].Address, uint64(int(BcLength)))
			}
		}

		LastBlock = block
		WriteBlock(BcLength, block)
	}
	IsMempAdded = false
	return isBlockValid
}

func GenMerkleRoot(transactions []transaction.Transaction) string {
	var height float64
	if len(transactions) == 1 {
		height = 1
	} else {
		height = math.Log2(float64(len(transactions))) + 1
		if float64(int(height)) != height {
			height = float64(int(height)) + 1
		}
	}

	var currLayer []string
	if len(transactions) != 0 {
		for i := 0; i < len(transactions); i++ {
			currLayer = append(currLayer, hashing.SHA1(transaction.GetCatTxFields(transactions[i])))
		}
	} else {
		currLayer = append(currLayer, hashing.SHA1(""))
	}

	for i := 0; i < int(height); i++ {
		var nextLayerLen int
		var isOddNodesCount bool = false
		if len(currLayer)%2 == 0 {
			nextLayerLen = len(currLayer) / 2
		} else {
			nextLayerLen = (len(currLayer) - 1) / 2
			isOddNodesCount = true
		}

		nextLayer := make([]string, nextLayerLen)

		currLayerInd := 0
		for j := 0; j < nextLayerLen; j++ {
			nextLayer[j] = hashing.SHA1(currLayer[currLayerInd] + currLayer[currLayerInd+1])
			currLayerInd += 2
		}

		if isOddNodesCount {
			nextLayer = append(nextLayer, hashing.SHA1(currLayer[len(currLayer)-1]))
		}
		currLayer = nextLayer
	}
	return currLayer[0]
}

func CreateBlock(rewardAdr byteArr.ByteArr, miningFlag int) Block {
	var newBlock Block

	if BcLength > 0 {
		newBlock.HashPrevBlock = hashing.SHA1(BlockToString(LastBlock))
	} else {
		newBlock.HashPrevBlock = "0000000000000000000000000000000000000000"
	}
	newBlock.Time = time.Now()
	var coinbaseTx transaction.Transaction
	coinbaseTx.Outputs = append(coinbaseTx.Outputs, transaction.Output{Address: rewardAdr, Value: GetCoinsForEmition()})

	var txArr []transaction.Transaction = GetTransactionsFromMempool(transaction.ComputeTxSize(coinbaseTx))

	var feeSum uint64 = 0
	for i := 0; i < len(txArr); i++ {
		feeSum += transaction.GetTxFee(txArr[i])
	}
	coinbaseTx.Outputs[0].Value += uint64(feeSum)
	IsMempAdded = true
	txArr = append([]transaction.Transaction{coinbaseTx}, txArr...)
	newBlock.Transactions = make([]transaction.Transaction, len(txArr))
	copy(newBlock.Transactions, txArr)

	newBlock.TransactionCount = uint(len(newBlock.Transactions))
	newBlock.MerkleRoot = GenMerkleRoot(newBlock.Transactions)
	log.Println("New block transaction count: " + fmt.Sprint(newBlock.TransactionCount))

	newBlock.Blocksize = uint(len(BlockToString(newBlock)))
	if miningFlag != -1 {
		newBlock.Bits = GetCurrBitsValue()
		log.Println("Current bits value is " + fmt.Sprintf("%x", newBlock.Bits))
		target := fmt.Sprintf("%x", BitsToTarget(newBlock.Bits))
		log.Println("Current target value is " + strings.Repeat("0", 40-len(target)) + target)
	}
	BlockForMining = newBlock
	if miningFlag == 0 {
		newBlock.Nonce = MineThreads(newBlock, 1)
	} else if miningFlag == 1 {
		newBlock.Nonce = MineThreads(newBlock, uint64(runtime.NumCPU()))
	}
	return newBlock
}

func ValidateBlock(block Block, height int) bool {
	var lastBlockHash string
	var prevBlock Block
	if height == int(BcLength) {
		prevBlock = LastBlock
	} else {
		var isBlockFound bool
		prevBlock, isBlockFound = GetBlock(uint64(height))
		if !isBlockFound {
			return false
		}
	}

	if int(BcLength) != 0 {
		lastBlockHash = hashing.SHA1(BlockToString(prevBlock))
	} else {
		lastBlockHash = "0000000000000000000000000000000000000000"
	}
	merkleRoot := GenMerkleRoot(block.Transactions)

	// Check bits value
	bits := GetCurrBitsValue()
	if bits != block.Bits {
		return false
	}

	// Check nonce
	hashNonce, _ := new(big.Int).SetString(hashing.SHA1(BlockToString(block)+fmt.Sprint(block.Nonce)), 16)
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
	for i := 1; i < int(block.TransactionCount); i++ {
		if !transaction.VerifyTransaction(block.Transactions[i]) {
			return false
		}
	}
	return true
}

func InitBlockchain() {
	BcLength, _ = GetBcHeight()
	if BcLength != 0 {
		LastBlock, _ = GetBlock(BcLength - 1)

		log.Println("Data is restored from db. Blockchain height:", BcLength)
	} else {
		FormGenesisBlock()
	}
}

func FormGenesisBlock() {
	log.Println("Creating initial block")

	var rewardAdr byteArr.ByteArr
	rewardAdr.SetFromHexString("e930fca003a4a70222d916a74cc851c3b3a9b050", 20)
	genesisBlock := CreateBlock(rewardAdr, 1)
	genesisBlock.Bits = STARTBITS

	if AddBlockToBlockchain(genesisBlock) {
		log.Println("Block is added to blockchain. Current height: " + fmt.Sprint(int(BcLength)+1) + "\n")
	} else {
		log.Println("Block is not added\n")
	}
	IncrBcHeight()
}

func PrintBlockTitle(block Block, height int) {
	println("Block height:", height)
	println("Version:", block.Version)
	println("Hash of prev block:", block.HashPrevBlock)
	println("Time:", block.Time.String())
	println("Merkle root:", block.MerkleRoot)
	println("Bits:", block.Bits)
	println("Nonce:", block.Nonce)
	println("Transaction count:", block.TransactionCount)
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

func (block *Block) SetFromByteArr(byteArr []byte) bool {
	var network bytes.Buffer
	dec := gob.NewDecoder(&network)
	err := dec.Decode(&block)
	if err != nil {
		log.Fatal("decode error:", err)
		return false
	}

	return true
}
