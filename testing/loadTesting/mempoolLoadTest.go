/*
	Tests how long does it wait for user to add tx to the mempool
	when there`s an incorrect tx spam
*/

package loadTesting

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/networking"
	"bavovnacoin/node/node_settings"
	"bavovnacoin/testing"
	"bavovnacoin/testing/account"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type MempoolLoadTest struct {
	testWorking bool

	txAmmount          uint
	incorrectTxAmmount uint

	currTxCounter int
	incTxCounter  int
	incTxId       float32
	prevIncTxId   float32

	txSendTime      []time.Time
	txAcceptionTime []time.Time

	source rand.Source
	random *rand.Rand
}

// TODO: send to a peer (p2p)
func (mlt *MempoolLoadTest) SendTxsToMempool() {
	var conn networking.Connection
	conn.Establish("localhost:25565")

	for uint(mlt.currTxCounter) < mlt.txAmmount {
		account.Wallet = append(account.Wallet, account.GenAccount(fmt.Sprint(mlt.currTxCounter)))
		testing.GenUtxo(mlt.currTxCounter, mlt.random)

		tx := testing.GenValidTx(mlt.currTxCounter, 1, mlt.random)

		mlt.currTxCounter++
		mlt.prevIncTxId = mlt.incTxId
		mlt.incTxId += float32(mlt.incorrectTxAmmount) / float32(mlt.txAmmount)

		if int(mlt.prevIncTxId) != int(mlt.incTxId) {
			tx, _ = testing.MakeTxIncorrect(tx, mlt.incTxCounter, mlt.random)
			mlt.incTxCounter++
		} else {
			mlt.txSendTime = append(mlt.txSendTime, time.Now())
		}

		var isAccepted bool
		conn.SendTransaction(tx, &isAccepted)
	}

	time.Sleep(100 * time.Millisecond)
	mlt.testWorking = false
}

func (mlt *MempoolLoadTest) countTxsData() {
	prevMempLen := 0
	for mlt.testWorking {
		if prevMempLen != len(blockchain.Mempool) {
			prevMempLen = len(blockchain.Mempool)
			mlt.txAcceptionTime = append(mlt.txAcceptionTime, time.Now())
		}
	}
}

func (mlt *MempoolLoadTest) printResult() {
	sectionsAmmount := 3
	if len(mlt.txAcceptionTime) < 3 {
		sectionsAmmount = len(mlt.txAcceptionTime)
	}

	var spentSections []int64
	var spentSectionsString []string
	var mean int64
	var ind int
	var valAmmount int64

	for i := 0; i < len(mlt.txAcceptionTime); i++ {
		mean += mlt.txAcceptionTime[i].Sub(mlt.txSendTime[i]).Milliseconds()
		valAmmount++

		if (i%len(mlt.txAcceptionTime)/sectionsAmmount == 0 || i == len(mlt.txAcceptionTime)-1) && i > 0 {
			time := mean / valAmmount
			spentSections = append(spentSections, time)
			spentSectionsString = append(spentSectionsString, fmt.Sprint(time))
			valAmmount = 0
			mean = 0
			ind++
		}
	}

	fmt.Printf("Mean time, divided by %d sections (ms): %s\n", sectionsAmmount, strings.Join(spentSectionsString, ", "))

	for i := 0; i < len(spentSections); i++ {
		mean += spentSections[i]
	}

	println("Mean time (ms):", mean/int64(len(spentSections)))
}

func (mlt *MempoolLoadTest) Launch() {
	println("Test started. Please, wait...")
	mlt.testWorking = true

	mlt.txAmmount = 100
	mlt.incorrectTxAmmount = 96

	testing.InitTestDb(false)
	node_settings.Settings.RPCip = "localhost:25565"
	networking.StartRPCListener()

	mlt.source = rand.NewSource(time.Now().Unix())
	mlt.random = rand.New(mlt.source)

	go mlt.SendTxsToMempool()
	mlt.countTxsData()

	println("Test ended. Results:")
	mlt.printResult()
}
