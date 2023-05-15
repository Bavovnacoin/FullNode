/*
	Checks txo and utxo values when confirming txs
*/

package singleFunctionTesting

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/dbController"
	"bavovnacoin/hashing"
	"bavovnacoin/networking_p2p"
	"bavovnacoin/node/node_controller/command_executor"
	"bavovnacoin/node/node_settings"
	"bavovnacoin/node/node_validator"
	"bavovnacoin/testing/account"
	"math/rand"
	"os"
	"time"
)

type TxoManagingTest struct {
	SingleFunctionTesting
	mcBlockAmmount uint64

	source rand.Source
	random *rand.Rand
}

func (tmt *TxoManagingTest) genBlocks() {
	command_executor.ComContr.FullNodeWorking = true
	go tmt.nodeWorkListener(tmt.mcBlockAmmount)
	go tmt.txForming("abc", tmt.random)
	node_validator.BlockGen(false)
}

func (tmt *TxoManagingTest) Launch() {
	tmt.mcBlockAmmount = 10

	node_settings.Settings.GetSettings()
	networking_p2p.Peer.StartP2PCommunication(node_settings.Settings.GetPrivKey(), node_settings.Settings.MyAddress, node_settings.Settings.OtherNodesAddresses)

	dbController.DbPath = "testing/testData"
	if _, err := os.Stat(dbController.DbPath); err == nil {
		os.RemoveAll(dbController.DbPath)
		println("Removed test db from a previous test.")
	}
	dbController.DB.OpenDb()

	blockchain.STARTBITS = 0xffff13
	tmt.source = rand.NewSource(time.Now().Unix())
	tmt.random = rand.New(tmt.source)

	tmt.genBlockTestAccounts(1)
	account.CurrAccount = account.Wallet[0]
	node_settings.Settings.RewardAddress = hashing.SHA1(account.CurrAccount.KeyPairList[0].PublKey)

	tmt.genBlocks()
	result := tmt.checkMcTxo()

	result.PrintTestOutput()
	result.PrintTestResult()
}
