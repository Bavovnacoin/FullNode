package loadtesting

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/dbController"
	"bavovnacoin/node"
	"bavovnacoin/testing"
	"bavovnacoin/transaction"
	"log"
	"math/rand"
	"os"
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

	testNodeWorking bool
	source          rand.Source
	random          *rand.Rand
}

func (lt *LoadTest) GenRandTx() transaction.Transaction {
	newTx := testing.GenValidTx(lt.currAccInd, 2, lt.random)

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
	for ; lt.txAmmount > 0; lt.txAmmount-- {
		newTx := lt.GenRandTx()
		conn.SendTransaction(newTx, &isAccepted)
	}
}

func (lt *LoadTest) StartLoadTest(txAmmount int, incTxAmmount int, rpcAmmount int) {
	lt.txAmmount = txAmmount
	lt.incTxAmmount = incTxAmmount
	lt.rpcAmmount = rpcAmmount
	lt.source = rand.NewSource(time.Now().Unix())
	lt.random = rand.New(lt.source)
	lt.testNodeWorking = true

	lt.initRandTxValues()

	testing.GenTestAccounts(txAmmount)
	testing.GenTestUtxo(txAmmount, lt.random)

	dbController.DbPath = "testing/loadTesting/testData"

	println(dbController.DB.OpenDb())
	node.StartRPC()

	go lt.StartTestTxSending()
	println("Initializing transactions. Please, wait...")
	time.Sleep(1 * time.Second)

	for lt.txAmmount != 0 || len(blockchain.Mempool) != 0 || len(blockchain.BlockForMining.Transactions) != 1 { //lt.testNodeWorking &&
		isAdded := node.AddBlock(false)
		if isAdded {
			println(len(blockchain.Mempool))
			log.Printf("Block is added to blockchain. Current height: %d. Handled %d test transactions\n", blockchain.BcLength, len(blockchain.LastBlock.Transactions)-1)
		}
	}

	dbController.DB.CloseDb()
	os.RemoveAll(dbController.DbPath)
	// go lt.TestNodeProcess()
	// node_controller.CommandHandler()
}
