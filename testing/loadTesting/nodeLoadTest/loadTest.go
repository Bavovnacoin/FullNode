/*
	Loads a network to check time that is needed for
	variety of node`s functions
*/

package nodeLoadTest

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/byteArr"
	"bavovnacoin/dbController"
	"bavovnacoin/hashing"
	"bavovnacoin/networking"
	"bavovnacoin/networking_p2p"
	"bavovnacoin/node/node_controller/command_executor"
	"bavovnacoin/node/node_validator"
	"bavovnacoin/testing"
	"bavovnacoin/testing/account"
	"bavovnacoin/txo"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type LoadTest struct {
	txAmmount  int
	rpcAmmount int

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
	rpcExecTimeUtxoByAddr  []time.Duration // Determines how much time a user needs to wait to call an rpc
	rpcExecTimeisAddrExist []time.Duration // Determines how much time a user needs to wait to call an rpc
	txVerifTimeP2P         []time.Duration // Determines how much time a user needs to wait to add tx to a mempool
	txVerifTimeRPC         []time.Duration // Determines how much time a user needs to wait to add tx to a mempool

	source rand.Source
	random *rand.Rand
}

func (lt *LoadTest) initTestData(txAmmount, rpcAmmount int) {
	lt.txAmmount = txAmmount
	lt.rpcAmmount = rpcAmmount
	lt.source = rand.NewSource(time.Now().Unix())
	lt.random = rand.New(lt.source)
	lt.rpcStep = (lt.txAmmount) / lt.rpcAmmount
}

func (lt *LoadTest) GenRandValues() {
	go testing.GenTestAccounts(lt.txAmmount)
	time.Sleep(5 * time.Second)
	go testing.GenTestUtxo(lt.txAmmount, lt.random)
}

func waitForMempoolChange() {
	prevAmmount := len(blockchain.Mempool)
	for prevAmmount == len(blockchain.Mempool) {
		time.Sleep(10 * time.Microsecond)
	}
}

func (lt *LoadTest) startTestTxSending() {
	var conn networking.Connection
	conn.Establish("localhost:8080")

	var newPeer networking_p2p.PeerData
	newPeer.Peer, _ = startTestPeer()
	newPeer, _ = addOtherAddress(fmt.Sprintf("%s/%s", networking_p2p.Peer.Peer.Addrs()[0], networking_p2p.Peer.Peer.ID().Pretty()), newPeer)

	var isAccepted bool
	var start time.Time
	txToCreate := lt.txAmmount
	for ; txToCreate > 0; txToCreate-- {
		println(txToCreate)
		newTx := testing.GenValidTx(lt.currAccInd, 2, lt.random)
		println("created")

		start = time.Now()
		if lt.currAccInd%2 == 0 {
			println("sending 1")
			conn.SendTransaction(newTx, &isAccepted)
			if isAccepted {
				lt.txVerifTimeRPC = append(lt.txVerifTimeRPC, time.Since(start))
				println("accepted")
			} else {
				println(blockchain.AreInputsInMempool(newTx.Inputs))
				println("not accepted")
			}
		} else {
			println("sending 2")
			newPeer.ProposeNewTx(newTx, "")
			//waitForMempoolChange()
			lt.txVerifTimeP2P = append(lt.txVerifTimeP2P, time.Since(start))
			println("accepted")
		}
		lt.currAccInd++
		println("----")
	}
}

func (lt *LoadTest) callRandRpc(rpcInd int) {
	var addr byteArr.ByteArr
	addr.SetFromHexString(hashing.SHA1("Glory to Ukraine"), 20)

	var conn Connection

	start := time.Now()
	if rpcInd%2 == 0 {
		conn.GetUtxoByAddress([]byteArr.ByteArr{addr})
		lt.rpcExecTimeUtxoByAddr = append(lt.rpcExecTimeUtxoByAddr, time.Since(start))
	}
	// else if rpcInd%2 == 1 {
	// 	conn.IsAddrExist(addr)
	// 	lt.rpcExecTimeisAddrExist = append(lt.rpcExecTimeisAddrExist, time.Since(start))
	// }
}

func (lt *LoadTest) tryCallRandRpc() {
	for ; lt.rpcHandledAmmount < lt.rpcAmmount; lt.rpcHandledAmmount++ {
		lt.callRandRpc(lt.rpcHandledAmmount)
		lt.rpcHandledAmmount++
	}

}

func (lt *LoadTest) testAddBlock() bool {
	if blockchain.AllowCreateBlock {
		var prevHash string
		if blockchain.BcLength > 0 {
			prevHash = hashing.SHA1(blockchain.BlockHeaderToString(blockchain.LastBlock))
		} else {
			prevHash = "0000000000000000000000000000000000000000"
		}

		node_validator.CreateBlockLog(blockchain.GetBits(false), prevHash, blockchain.LastBlock, false)
		blockchain.AllowCreateBlock = false
	}

	if blockchain.CreatedBlock.MerkleRoot != "" { // Is block mined check
		isBlockValid := blockchain.VerifyBlock(blockchain.CreatedBlock, int(blockchain.BcLength), true, false)
		isAdded := node_validator.AddBlockLog(false, isBlockValid)
		blockchain.CreatedBlock.MerkleRoot = ""
		return isAdded
	}
	command_executor.PauseCommand()
	return false
}

func (lt *LoadTest) Launch(txAmmount int, rpcAmmount int) {
	lt.initTestData(txAmmount, rpcAmmount)
	InitTestSettings()

	dbController.DbPath = "testing/testData"
	if _, err := os.Stat(dbController.DbPath); err == nil {
		os.RemoveAll(dbController.DbPath)
		println("Removed test db from a previous test.")
	}
	dbController.DB.OpenDb()

	networking_p2p.Peer.StartP2PCommunication()
	node_validator.StartRPC()

	testing.GenTestAccounts(lt.txAmmount)
	fmt.Printf("Generated %d test accounts\n", len(account.Wallet))
	testing.GenTestUtxo(lt.txAmmount, lt.random)
	fmt.Printf("Generated %d test utxo\n", len(txo.CoinDatabase))

	lt.startTestTxSending() // add go
	//go lt.tryCallRandRpc()
	println("Started tx sending")

	// println("Generating blocks")
	// var handledTxAmmount int
	// var isAdded bool
	// for (lt.txAmmount != 0 || len(blockchain.Mempool) != 0 || len(blockchain.BlockForMining.Transactions) != 0) && (handledTxAmmount < lt.txAmmount) { //|| lt.rpcHandledAmmount < lt.rpcAmmount
	// 	isAdded = lt.testAddBlock()

	// 	if isAdded {
	// 		lt.txHandled += len(blockchain.LastBlock.Transactions) - 1
	// 		log.Printf("Block is added to blockchain. Current height: %d. Handled %d transactions\n",
	// 			blockchain.BcLength, len(blockchain.LastBlock.Transactions))

	// 		handledTxAmmount += len(blockchain.LastBlock.Transactions)
	// 		blockchain.BlockForMining = blockchain.Block{}
	// 	} else {
	// 		println("Block is not added")
	// 	}

	// }
	//println(handledTxAmmount, lt.txAmmount, handledTxAmmount < lt.txAmmount)

	println(len(blockchain.Mempool))

	dbController.DB.CloseDb()
	os.RemoveAll(dbController.DbPath)

	lt.printResults()
}
