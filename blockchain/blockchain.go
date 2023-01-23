package blockchain

import (
	"bavovnacoin/hashing"
	"bavovnacoin/transaction"
	"bavovnacoin/utxo"
	"fmt"
	"math"
	"math/big"
	"time"
)

var Blockchain []Block

type Block struct {
	Id               uint
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
	str += fmt.Sprint(block.Id)
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
	isBlockValid := ValidateBlock(block)
	if isBlockValid {
		for i := 0; i < len(block.Transactions); i++ {
			txInpList := block.Transactions[i].Inputs

			for j := 0; j < len(txInpList); j++ {
				utxo.DelFromUtxo(txInpList[j].HashAdr, txInpList[j].OutInd)
			}

			txOutList := block.Transactions[i].Outputs
			for j := 0; j < len(txOutList); j++ {
				utxo.AddToUtxo(txOutList[j].HashAdr, txOutList[j].Sum)
			}
		}
		Blockchain = append(Blockchain, block)
	}
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

func CreateBlock(id uint, txArr []transaction.Transaction, rewardAdr string, miningFlag int) Block {
	var newBlock Block
	newBlock.Id = id

	if len(Blockchain) > 0 {
		newBlock.HashPrevBlock = hashing.SHA1(BlockToString(Blockchain[len(Blockchain)-1]))
		newBlock.Time = time.Now()
	} else {
		newBlock.Time = time.Date(2022, 12, 31, 11, 54, 22, 0, time.Local)
		newBlock.HashPrevBlock = "0000000000000000000000000000000000000000"
	}

	newBlock.Transactions = make([]transaction.Transaction, len(txArr)+1)

	var coinbaseTx transaction.Transaction
	coinbaseTx.Outputs = append(coinbaseTx.Outputs, transaction.Output{HashAdr: rewardAdr, Sum: GetCoinsForEmition()})

	txArr = append([]transaction.Transaction{coinbaseTx}, txArr...)
	copy(newBlock.Transactions, txArr)

	newBlock.TransactionCount = uint(len(newBlock.Transactions))
	newBlock.MerkleRoot = GenMerkleRoot(newBlock.Transactions)

	newBlock.Blocksize = uint(len(BlockToString(newBlock)))

	bits := GetCurrBitsValue()

	if miningFlag == 0 {
		newBlock.Nonce = MineCommon(newBlock, bits)
	} else {
		newBlock.Nonce = MineAllThreads(newBlock, bits)
	}
	return newBlock
}

func ValidateBlock(block Block) bool {
	lastBlockHash := hashing.SHA1(BlockToString(Blockchain[len(Blockchain)-1]))
	merkleRoot := GenMerkleRoot(block.Transactions)
	// Transaction are verified when added to mempool, no need to double check it

	bits := GetCurrBitsValue()
	if bits != block.Bits {
		return false
	}

	hashNounce, _ := new(big.Int).SetString(hashing.SHA1(BlockToString(block)+fmt.Sprint(block.Nonce)), 16)
	if BitsToTarget(block.Bits).Cmp(hashNounce) != 1 {
		return false
	}

	if block.HashPrevBlock != lastBlockHash ||
		block.MerkleRoot != merkleRoot {
		return false
	}

	return true
}

func InitBlockchain() {
	genesisBlock := CreateBlock(0,
		nil, "0284cbd0bcf8a34035b71c5a72e37924cb960aaa0b69df4c41d50628734b8e1408", 1)

	Blockchain = append(Blockchain, genesisBlock)
}
