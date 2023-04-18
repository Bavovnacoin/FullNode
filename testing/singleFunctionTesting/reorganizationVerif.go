package singleFunctionTesting

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/dbController"
	"bavovnacoin/node"
	"bavovnacoin/node_controller/command_executor"
	"math/rand"
	"os"
	"time"
)

type ReorganizationVerifTest struct {
	SingleFunctionTesting
	blockAmmount int

	source rand.Source
	random *rand.Rand
}

func (rv *ReorganizationVerifTest) nodeWorkListener() {
	for true {
		if blockchain.BcLength >= uint64(rv.blockAmmount) {
			command_executor.ComContr.FullNodeWorking = false
			return
		}
	}
}

func (rv *ReorganizationVerifTest) genBlocks() {
	command_executor.ComContr.FullNodeWorking = true
	go rv.nodeWorkListener()
	node.BlockGen(false)
}

func (rv *ReorganizationVerifTest) genAltchBlocks() {

}

func (rv *ReorganizationVerifTest) Launch() {
	rv.blockAmmount = 3

	dbController.DbPath = "testing/testData"
	if _, err := os.Stat(dbController.DbPath); err == nil {
		os.RemoveAll(dbController.DbPath)
		println("Removed test db from a previous test.")
	}
	dbController.DB.OpenDb()

	blockchain.STARTBITS = 0xffff14
	rv.source = rand.NewSource(time.Now().Unix())
	rv.random = rand.New(rv.source)

	rv.genBlockTestAccounts(rv.blockAmmount)
	rv.genBlocks() // Generating blocks in mainchain

	rv.genAltchBlocks()
}
