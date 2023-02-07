package blockchain

import (
	"bavovnacoin/byteArr"
	"bavovnacoin/hashing"
	"bavovnacoin/transaction"
	"bavovnacoin/utxo"
	"fmt"
	"log"
	"math"
	"math/big"
	"runtime"
	"time"
)

var Blockchain []Block

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
	isBlockValid := ValidateBlock(block, len(Blockchain))

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
				utxo.AddUtxo(txByteArr, uint64(j), txOutList[j].Sum, txOutList[j].Address, uint64(len(Blockchain)))
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

func CreateBlock(rewardAdr byteArr.ByteArr, miningFlag int) Block {
	var newBlock Block

	if len(Blockchain) > 0 {
		newBlock.HashPrevBlock = hashing.SHA1(BlockToString(Blockchain[len(Blockchain)-1]))
	} else {
		newBlock.HashPrevBlock = "0000000000000000000000000000000000000000"
	}
	newBlock.Time = time.Now()
	var coinbaseTx transaction.Transaction
	coinbaseTx.Outputs = append(coinbaseTx.Outputs, transaction.Output{Address: rewardAdr, Sum: GetCoinsForEmition()})

	var txArr []transaction.Transaction = GetTransactionsFromMempool(transaction.ComputeTxSize(coinbaseTx))

	var feeSum uint64 = 0
	for i := 0; i < len(txArr); i++ {
		feeSum += transaction.GetTxFee(txArr[i])
	}
	coinbaseTx.Outputs[0].Sum += uint64(feeSum)

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
	}

	if miningFlag == 0 {
		newBlock.Nonce = MineThreads(newBlock, 1)
	} else if miningFlag == 1 {
		newBlock.Nonce = MineThreads(newBlock, uint64(runtime.NumCPU()))
	}
	return newBlock
}

func ValidateBlock(block Block, id int) bool {
	var lastBlockHash string
	if len(Blockchain) != 0 {
		lastBlockHash = hashing.SHA1(BlockToString(Blockchain[len(Blockchain)-1]))
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
	if !CheckEmitedCoins(block.Transactions[0].Outputs[0].Sum-allFee, id) {
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
	log.Println("Creating initial block")

	var rewardAdr byteArr.ByteArr
	rewardAdr.SetFromHexString("e930fca003a4a70222d916a74cc851c3b3a9b050", 20)
	genesisBlock := CreateBlock(rewardAdr, 1)
	genesisBlock.Bits = STARTBITS

	if AddBlockToBlockchain(genesisBlock) {
		log.Println("Block is added to blockchain. Current length: " + fmt.Sprint(len(Blockchain)) + "\n")
	} else {
		log.Println("Block is not added\n")
	}
}

func PrintBlockTitle(block Block, id int) {
	println("Block id:", id)
	println("Version:", block.Version)
	println("Hash of prev block:", block.HashPrevBlock)
	println("Time:", block.Time.String())
	println("Merkle root:", block.MerkleRoot)
	println("Bits:", block.Bits)
	println("Nonce:", block.Nonce)
	println("Transaction count:", block.TransactionCount)
}
