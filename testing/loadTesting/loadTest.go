package loadtesting

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/byteArr"
	"bavovnacoin/dbController"
	"bavovnacoin/hashing"
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
	currAccInd   int
	txHandled    int

	// Rand rpc generation
	rpcStep                 int
	rpcLastTxHandledAmmount int
	rpcHandledAmmount       int

	// Test results
	rpcExecTime []time.Duration
	txVerifTime []time.Duration // Determines how much time a user needs to wait to add tx to a mempool

	source rand.Source
	random *rand.Rand
}

func (lt *LoadTest) initTestData(txAmmount, incTxAmmount, rpcAmmount int) {
	lt.txAmmount = txAmmount
	lt.incTxAmmount = incTxAmmount
	lt.rpcAmmount = rpcAmmount
	lt.source = rand.NewSource(time.Now().Unix())
	lt.random = rand.New(lt.source)
	lt.rpcStep = (lt.txAmmount - lt.incTxAmmount) / lt.rpcAmmount
	lt.initRandTxValues()
}

func (lt *LoadTest) genRandTx() transaction.Transaction {
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

func (lt *LoadTest) startTestTxSending() {
	var conn Connection
	conn.Establish()
	defer conn.Close()

	var isAccepted bool
	var start time.Time
	for ; lt.txAmmount > 0; lt.txAmmount-- {
		newTx := lt.genRandTx()

		start = time.Now()
		conn.SendTransaction(newTx, &isAccepted)
		lt.txVerifTime = append(lt.txVerifTime, time.Since(start))

	}
}

func (lt *LoadTest) callRandRpc(rpcInd int) {
	var addr byteArr.ByteArr
	addr.SetFromHexString(hashing.SHA1("Glory to Ukraine"), 20)

	var conn Connection
	conn.Establish()
	defer conn.Close()

	start := time.Now()
	if rpcInd%2 == 0 {
		conn.GetUtxoByAddress([]byteArr.ByteArr{addr})
	} else if rpcInd%2 == 1 {
		conn.IsAddrExist(addr)
	}
	lt.rpcExecTime = append(lt.rpcExecTime, time.Since(start))
}

func (lt *LoadTest) tryCallRandRpc() {
	currCallsAmmount := (lt.txHandled - lt.rpcStep*lt.rpcHandledAmmount) / lt.rpcStep

	if lt.rpcHandledAmmount < lt.rpcAmmount && currCallsAmmount != 0 {
		for i := 0; i < currCallsAmmount; i++ {
			lt.callRandRpc(lt.rpcHandledAmmount)
			lt.rpcHandledAmmount++
		}
	}

}

func (lt *LoadTest) printResults() {

}

func (lt *LoadTest) StartLoadTest(txAmmount int, incTxAmmount int, rpcAmmount int) {
	lt.initTestData(txAmmount, incTxAmmount, rpcAmmount)

	dbController.DbPath = "testing/loadTesting/testData"
	if _, err := os.Stat(dbController.DbPath); err == nil {
		os.RemoveAll(dbController.DbPath)
		println("Removed test db from a previous test.")
	}
	dbController.DB.OpenDb()

	testing.GenTestAccounts(txAmmount)
	testing.GenTestUtxo(txAmmount, lt.random)

	node.StartRPC()

	go lt.startTestTxSending()
	println("Initializing transactions. Please, wait...")
	time.Sleep(1 * time.Second)

	// Add rpc check (is it executing -> store ammount of executed time as an array and compare to ammount of need rpc?)
	for lt.txAmmount != 0 || len(blockchain.Mempool) != 0 || len(blockchain.BlockForMining.Transactions) != 1 ||
		len(lt.rpcExecTime) < lt.rpcAmmount {
		isAdded := node.AddBlock(false)
		if isAdded {
			lt.txHandled += len(blockchain.LastBlock.Transactions) - 1
			log.Printf("Block is added to blockchain. Current height: %d. Handled %d test transactions\n",
				blockchain.BcLength, len(blockchain.LastBlock.Transactions)-1)
		}

		lt.tryCallRandRpc()
	}

	dbController.DB.CloseDb()
	os.RemoveAll(dbController.DbPath)

	lt.printResults()
	// node_controller.CommandHandler()
}
