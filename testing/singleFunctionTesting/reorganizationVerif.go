package singleFunctionTesting

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/byteArr"
	"bavovnacoin/dbController"
	"bavovnacoin/hashing"
	"bavovnacoin/networking_p2p"
	"bavovnacoin/node/node_controller/command_executor"
	"bavovnacoin/node/node_controller/node_settings"
	"bavovnacoin/node/node_validator"
	"bavovnacoin/testing/account"
	"bavovnacoin/transaction"
	"bavovnacoin/txo"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type ReorganizationVerifTest struct {
	SingleFunctionTesting
	mcBlockAmmount uint64
	acBlockAmmount uint64
	prevHeight     uint64
	currAccKpInd   int

	source rand.Source
	random *rand.Rand
}

func (rv *ReorganizationVerifTest) CreateTx() transaction.Transaction {
	fee := rv.random.Intn(5) + 1
	isGenLocktime := rv.random.Intn(5)
	var locktime uint
	if isGenLocktime == 2 {
		locktime = uint(int(blockchain.BcLength+1) + rv.random.Intn(2) + 1)
	}

	var outAddr []byteArr.ByteArr
	outAddr = append(outAddr, byteArr.ByteArr{})
	outAddr[0].SetFromHexString(hashing.SHA1("An address"), 20) // account.Wallet[0].KeyPairList[rv.currAccKpInd].PublKey

	var outValue []uint64
	outValue = append(outValue, 1000)

	tx, _ := transaction.CreateTransaction(fmt.Sprint(0), outAddr, outValue, fee, locktime)
	rv.currAccKpInd++
	return tx
}

func (rv *ReorganizationVerifTest) nodeWorkListener(blocksCount uint64) {
	for true {
		if blockchain.BcLength >= blocksCount {
			command_executor.ComContr.FullNodeWorking = false
			return
		}

		if rv.prevHeight != blockchain.BcLength {
			rv.prevHeight = blockchain.BcLength
			blockchain.AddTxToMempool(rv.CreateTx(), false)
		}
	}
}

func (rv *ReorganizationVerifTest) genBlocks() {
	command_executor.ComContr.FullNodeWorking = true
	go rv.nodeWorkListener(rv.mcBlockAmmount)
	node_validator.BlockGen(false)
}

func (rv *ReorganizationVerifTest) genAltchBlocks() {
	bl, _ := blockchain.GetBlock(blockchain.BcLength-2, 0)
	var prevHash string = hashing.SHA1(blockchain.BlockHeaderToString(bl))

	for i := uint64(0); i < rv.acBlockAmmount; i++ {
		blockchain.AddTxToMempool(rv.CreateTx(), false)
		node_validator.CreateBlockLog(blockchain.GetBits(true), prevHash, bl, true)
		blockchain.AllowCreateBlock = false

		var otherNodesTime []int64
		otherNodesTime = append(otherNodesTime, time.Now().UTC().Unix())
		if i == 0 {
			blockchain.CreatedBlock.Time = bl.Time
		}
		blockchain.CreatedBlock.Version = 1
		blockchain.TryCameBlockToAdd(blockchain.CreatedBlock, blockchain.BcLength-1+uint64(i), otherNodesTime)
		prevHash = hashing.SHA1(blockchain.BlockHeaderToString(blockchain.CreatedBlock))
		bl = blockchain.CreatedBlock
	}
}

func getOutputsFromMainchain() ([]txo.TXO, []txo.TXO, bool) {
	var mainchTxo []txo.TXO
	var mainchUtxo []txo.TXO
	for i := uint64(0); i < blockchain.BcLength; i++ {
		block, res := blockchain.GetBlock(i, 0)
		if !res {
			print("Test failed. Error when checking mainchain")
			return mainchTxo, mainchUtxo, false
		}

		for j := 0; j < len(block.Transactions); j++ {
			// Check inputs
			inputs := block.Transactions[j].Inputs
			for k := 0; k < len(inputs); k++ {
				for h := 0; h < len(mainchUtxo); h++ {
					if inputs[k].TxHash.IsEqual(mainchUtxo[h].OutTxHash) &&
						inputs[k].OutInd == int(mainchUtxo[h].TxOutInd) {
						mainchTxo = append(mainchTxo, txo.TXO{OutTxHash: mainchUtxo[h].OutTxHash, TxOutInd: mainchUtxo[h].TxOutInd})
						mainchUtxo = append(mainchUtxo[:h], mainchUtxo[h+1:]...)
						break
					}
				}

			}

			// Check outputs
			outputs := block.Transactions[j].Outputs
			for k := 0; k < len(outputs); k++ {
				mainchUtxo = append(mainchUtxo, txo.TXO{Value: outputs[k].Value, OutAddress: outputs[k].Address, BlockHeight: i})
			}
		}
	}

	return mainchTxo, mainchUtxo, true
}

func (rv *ReorganizationVerifTest) printResult() {
	println("Results:")
	println("Blockchain scheme:")
	for height := 0; true; height++ {
		blocks, res := blockchain.GetBlocksOnHeight(uint64(height))
		if !res || len(blocks) == 0 {
			break
		}

		var str string
		if len(blocks) == 1 {
			if blocks[0].ChainId == 0 {
				str += fmt.Sprint(blocks[0].Block.Version)
				str += "  "
			} else {
				str += "  "
				str += fmt.Sprint(blocks[0].Block.Version)
			}
		} else {
			str += fmt.Sprintf("%s %s", fmt.Sprint(blocks[0].Block.Version), fmt.Sprint(blocks[1].Block.Version))
		}

		// println(str)
	}

	// Check TXO and UTXO
	mTxo, mUtxo, res := getOutputsFromMainchain()
	if !res {
		return
	}

	storedTxo, _ := txo.GetTxoList("txo")
	storedUtxo, _ := txo.GetTxoList("utxo")
	for i := 0; i < len(mUtxo); i++ {
		mUtxo[i].PrintTxo(-1)
	}
	println()
	for i := 0; i < len(storedUtxo); i++ {
		storedUtxo[i].PrintTxo(-1)
	}

	println(len(mUtxo), len(storedUtxo), len(txo.CoinDatabase))
	if len(mTxo) != len(storedTxo) {
		println("Test failed. Wrong txo ammount.")
		return
	} else if len(mUtxo) != len(storedUtxo) {
		println("Test failed. Wrong utxo ammount.")
		//return
	}

	for i := 0; i < len(storedTxo); i++ {
		_, res := txo.GetTxo(storedTxo[i].OutTxHash, int(storedTxo[i].TxOutInd), storedTxo[i].BlockHeight)
		if !res {
			println("Incorrect txo")
			return
		}
	}

	for i := 0; i < len(storedUtxo); i++ {
		_, res := txo.GetUtxos(storedUtxo[i].OutTxHash, int(storedUtxo[i].TxOutInd))
		if !res {
			println("Incorrect utxo")
			return
		}
		// for j := 0; j < len(utxos); j++ {
		// 	utxos[j].PrintTxo(i)
		// }
	}
	println("Test passed!")
}

func (rv *ReorganizationVerifTest) Launch() {
	rv.mcBlockAmmount = 4
	rv.acBlockAmmount = 3
	node_settings.Settings.GetSettings()
	networking_p2p.StartP2PCommunication()

	dbController.DbPath = "testing/testData"
	if _, err := os.Stat(dbController.DbPath); err == nil {
		os.RemoveAll(dbController.DbPath)
		println("Removed test db from a previous test.")
	}
	dbController.DB.OpenDb()

	blockchain.STARTBITS = 0xffff14
	rv.source = rand.NewSource(time.Now().Unix())
	rv.random = rand.New(rv.source)

	rv.genBlockTestAccounts(1) //int(rv.mcBlockAmmount) + int(rv.acBlockAmmount)
	account.CurrAccount = account.Wallet[0]
	node_settings.Settings.RewardAddress = hashing.SHA1(account.CurrAccount.KeyPairList[0].PublKey)

	blockchain.STARTBITS = 0xffff13
	rv.genBlocks() // Generating blocks in mainchain

	// rv.genAltchBlocks()
	time.Sleep(time.Millisecond * 800)
	rv.printResult()
}
