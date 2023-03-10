package loadtesting

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/node"
	"bavovnacoin/node_controller/command_executor"
	"bavovnacoin/testing"
	"bavovnacoin/transaction"
	"log"
	"math/rand"
	"time"
)

type LoadTest struct {
	txAmmount    int
	incTxAmmount int
	rpcAmmount   int

	// Rand tx generation
	step         int
	incTxInd     int
	incTxCounter int
	currAccInd   int // TODO: change to currAccInd

	nodeWorking bool
	source      rand.Source
	random      *rand.Rand
	blockchain  []blockchain.Block
}

func (lt *LoadTest) GenRandTx() transaction.Transaction {
	newTx := testing.GenValidTx(lt.currAccInd, lt.random)

	if lt.currAccInd == lt.incTxInd && lt.incTxCounter <= lt.incTxAmmount {
		stStep := lt.step * lt.incTxCounter
		if lt.incTxAmmount-1 == lt.incTxCounter {
			lt.step = lt.txAmmount - stStep
		}
		lt.incTxInd = lt.random.Intn(lt.step) + stStep

		newTx, _ = testing.MakeTxIncorrect(newTx, lt.incTxCounter, lt.random)
		lt.incTxCounter++
	}

	lt.currAccInd++

	return newTx
}

func (lt *LoadTest) initRandTxValues() {
	lt.step = int(lt.txAmmount / lt.incTxAmmount)
	lt.incTxInd = -1
	if lt.incTxAmmount != 0 {
		stStep := lt.step * lt.incTxCounter
		lt.incTxInd = lt.random.Intn(lt.step) + stStep

		lt.incTxCounter++
	}
}

func (lt *LoadTest) StartTestTxSending() {
	var conn Connection
	conn.Establish()
	defer conn.Close()

	var isAccepted bool
	for i := 0; i < lt.txAmmount; i++ {
		newTx := lt.GenRandTx()
		conn.SendTransaction(newTx, &isAccepted)
	}
}

func (lt *LoadTest) GetCurrBitsValue() uint64 {
	var bits uint64
	if int(blockchain.BcLength)%blockchain.BLOCK_DIFF_CHECK == 0 && blockchain.BcLength != 0 {
		blockDiff := lt.blockchain[uint64(int(blockchain.BcLength)-blockchain.BLOCK_DIFF_CHECK)]
		bits = blockchain.GenBits(blockDiff.Time, blockchain.LastBlock.Time, blockchain.LastBlock.Bits)
	} else if blockchain.BcLength != 0 {
		bits = blockchain.LastBlock.Bits
	} else {
		bits = blockchain.STARTBITS
	}
	return bits
}

func (lt *LoadTest) StartLoadTest(txAmmount int, incTxAmmount int, rpcAmmount int) {
	lt.txAmmount = txAmmount
	lt.incTxAmmount = incTxAmmount
	lt.rpcAmmount = rpcAmmount
	lt.source = rand.NewSource(time.Now().Unix())
	lt.random = rand.New(lt.source)
	lt.nodeWorking = true

	lt.initRandTxValues()

	testing.GenTestAccounts(txAmmount)
	testing.GenTestUtxo(10, lt.random)

	node.StartRPC()

	go lt.StartTestTxSending()
	println("Initializing transactions. Please, wait...")
	time.Sleep(1 * time.Second)

	for lt.nodeWorking {
		if node.AllowCreateBlock {
			log.Println("Creating a new block")
			go node.CreateBlockLog(lt.GetCurrBitsValue())
			println(len(blockchain.Mempool))
			node.AllowCreateBlock = false
		}

		if node.CreatedBlock.MerkleRoot != "" { // Is block mined check
			node.CreatedBlock.Bits = lt.GetCurrBitsValue()
			isCreated := node.AddBlockLog(false)
			if isCreated {
				lt.blockchain = append(lt.blockchain, node.CreatedBlock)
			}
			node.CreatedBlock.MerkleRoot = ""
		}
		command_executor.PauseCommand()
	}

	// go lt.TestNodeProcess()
	// node_controller.CommandHandler()
}
