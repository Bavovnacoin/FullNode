package blockchain

import (
	"bavovnacoin/byteArr"
	"bavovnacoin/hashing"
	"bavovnacoin/transaction"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"math/big"
	"time"
)

type Block struct {
	Blocksize     uint
	Version       uint
	HashPrevBlock string
	Time          int64
	MerkleRoot    string
	Bits          uint64
	Nonce         uint64
	Chainwork     *big.Int
	Transactions  []transaction.Transaction
}

// TODO: add Time to header
func BlockHeaderToString(block Block) string {
	str := ""
	str += fmt.Sprint(block.Blocksize)
	str += fmt.Sprint(block.Version)
	str += block.HashPrevBlock
	str += block.MerkleRoot
	str += fmt.Sprintf("%x", block.Bits)
	str += fmt.Sprint(block.Nonce)
	str += block.Chainwork.String()
	return str
}

func getCurrBlockChainwork(block Block) *big.Int {
	blockTarget := BitsToTarget(block.Bits)
	return new(big.Int).Div(hashing.MaxNum, new(big.Int).Add(blockTarget, big.NewInt(1)))
}

func CreateBlock(rewardAdr byteArr.ByteArr, allowPrint bool) Block {
	var newBlock Block

	if BcLength > 0 {
		newBlock.HashPrevBlock = hashing.SHA1(BlockHeaderToString(LastBlock))
	} else {
		newBlock.HashPrevBlock = "0000000000000000000000000000000000000000"
	}

	var coinbaseTx transaction.Transaction
	coinbaseTx.Outputs = append(coinbaseTx.Outputs, transaction.Output{Address: rewardAdr, Value: GetCoinsForEmition()})

	txArr := GetTransactionsFromMempool(transaction.ComputeTxSize(coinbaseTx))

	var feeSum uint64 = 0
	for i := 0; i < len(txArr); i++ {
		feeSum += transaction.GetTxFee(txArr[i])
	}
	coinbaseTx.Outputs[0].Value += uint64(feeSum)
	IsMempAdded = true
	txArr = append([]transaction.Transaction{coinbaseTx}, txArr...)
	newBlock.Transactions = make([]transaction.Transaction, len(txArr))
	copy(newBlock.Transactions, txArr)

	newBlock.MerkleRoot = GenMerkleRoot(newBlock.Transactions)
	if allowPrint {
		log.Println("New block transaction count: " + fmt.Sprint(len(newBlock.Transactions)))
	}

	newBlock.Blocksize = uint(len(BlockHeaderToString(newBlock)))
	newBlock.Chainwork = getChainwork(newBlock)
	return newBlock
}

func PrintBlockTitle(block Block, height uint64) {
	println("Block height:", height)
	println("Version:", block.Version)
	println("Hash of prev block:", block.HashPrevBlock)
	println("Time:", time.Unix(block.Time, 0).String())
	println("Merkle root:", block.MerkleRoot)
	println("Bits:", block.Bits)
	println("Nonce:", block.Nonce)
	println("Transactions count:", len(block.Transactions))
	println("Chainwork:", block.Chainwork.String())
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
