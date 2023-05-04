/*
	Checks outputs and blockchain structure after reorganization
*/

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
	"strings"
	"time"
)

type ReorganizationVerifTest struct {
	SingleFunctionTesting
	mcBlockAmmount uint64
	acBlockAmmount uint64
	prevHeight     uint64

	source rand.Source
	random *rand.Rand
}

func (rv *ReorganizationVerifTest) txForming() {
	for command_executor.ComContr.FullNodeWorking {
		rv.prevHeight = blockchain.BcLength
		tx, isValid := rv.CreateTestTx("abc", rv.random)
		if !isValid {
			continue
		}
		blockchain.AddTxToMempool(tx, true)
	}
}

func (rv *ReorganizationVerifTest) nodeWorkListener(blocksCount uint64) {
	for true {
		if blockchain.BcLength >= blocksCount {
			command_executor.ComContr.FullNodeWorking = false
			return
		}
	}
}

func (rv *ReorganizationVerifTest) genBlocks() {
	command_executor.ComContr.FullNodeWorking = true
	go rv.nodeWorkListener(rv.mcBlockAmmount)
	go rv.txForming()
	node_validator.BlockGen(false)
}

func CreateBlock(bits uint64, prevHash string, lastBlock blockchain.Block, allowPrint bool) {
	var rewardAdr byteArr.ByteArr
	rewardAdr.SetFromHexString(node_settings.Settings.RewardAddress, 20)
	newBlock := blockchain.CreateBlock(rewardAdr, prevHash, allowPrint)
	newBlock.Bits = bits
	newBlock.Chainwork = blockchain.GetChainwork(newBlock, lastBlock)
	var miningRes bool

	newBlock, miningRes = blockchain.MineBlock(newBlock, 1, allowPrint)

	if !miningRes {
		blockchain.AllowCreateBlock = true
		return
	}

	blockchain.IsMiningDone = true
	blockchain.RemoveTxsFromMempool(newBlock.Transactions[1:])
	blockchain.CreatedBlock = newBlock
}

func (rv *ReorganizationVerifTest) genAltchBlocks() {
	bl, _ := blockchain.GetBlock(blockchain.BcLength-2, 0)
	var prevHash string = hashing.SHA1(blockchain.BlockHeaderToString(bl))

	for i := uint64(0); i < rv.acBlockAmmount; i++ {
		node_settings.Settings.RewardAddress = hashing.SHA1("abc")
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
			var txByteArr byteArr.ByteArr
			txByteArr.SetFromHexString(hashing.SHA1(transaction.GetCatTxFields(block.Transactions[j])), 20)

			// Check inputs
			inputs := block.Transactions[j].Inputs
			for k := 0; k < len(inputs); k++ {
				for h := 0; h < len(mainchUtxo); h++ {
					if inputs[k].TxHash.IsEqual(mainchUtxo[h].OutTxHash) &&
						inputs[k].OutInd == int(mainchUtxo[h].TxOutInd) {
						mainchTxo = append(mainchTxo, mainchUtxo[h])
						println("added txo")
						mainchUtxo = append(mainchUtxo[:h], mainchUtxo[h+1:]...)
						break
					}
				}

			}

			// Check outputs
			outputs := block.Transactions[j].Outputs
			for k := 0; k < len(outputs); k++ {
				mainchUtxo = append(mainchUtxo, txo.TXO{Value: outputs[k].Value, OutAddress: outputs[k].Address,
					BlockHeight: i, OutTxHash: txByteArr, TxOutInd: uint64(k)})
			}
		}
	}

	return mainchTxo, mainchUtxo, true
}

func (rv *ReorganizationVerifTest) isBlocksChanged(levels [][]string) bool {
	if uint64(len(levels)) < rv.mcBlockAmmount+rv.acBlockAmmount {
		return false
	}

	for i := 0; i < len(levels); i++ {
		if len(levels[i]) == 1 {
			if i < int(rv.mcBlockAmmount) && levels[i][0] != "0" ||
				i > int(rv.mcBlockAmmount) && levels[i][0] != "1" {
				return false
			}
		} else {
			if levels[i][0] != "1" && levels[i][1] != "0" {
				return false
			}
		}
	}

	return true
}

func (rv *ReorganizationVerifTest) printResult() {
	println("Results:")
	println("Blockchain scheme:")
	var levels [][]string
	for height := 0; true; height++ {
		blocks, res := blockchain.GetBlocksOnHeight(uint64(height))
		if !res || len(blocks) == 0 {
			break
		}

		var level []string
		if len(blocks) == 1 {
			if blocks[0].ChainId == 0 {
				level = append(level, fmt.Sprint(blocks[0].Block.Version))
				level = append(level, " ")
			} else {
				level = append(level, " ")
				level = append(level, fmt.Sprint(blocks[0].Block.Version))
			}
		} else {
			level = append(level, fmt.Sprint(blocks[0].Block.Version))
			level = append(level, fmt.Sprint(blocks[1].Block.Version))
		}
		levels = append(levels, level)
		println(strings.Join(level, " "))
	}

	println("")
	if !rv.isBlocksChanged(levels) {
		println("The scheme is correct")
	} else {
		println("Test failed: the scheme is incorrect")
		return
	}

	// Check TXO and UTXO
	mTxo, mUtxo, res := getOutputsFromMainchain()
	if !res {
		return
	}

	storedTxo, _ := txo.GetTxoList("txo")
	storedUtxo, _ := txo.GetTxoList("utxo")

	if len(mTxo) != len(storedTxo) {
		println(len(mTxo), len(storedTxo))
		println("Test failed. Wrong txo ammount.")
		return
	} else if len(mUtxo) != len(storedUtxo) {
		println("Test failed. Wrong utxo ammount.")
		return
	}

	for i := 0; i < len(storedTxo); i++ {
		_, res := txo.GetTxos(storedTxo[i].OutTxHash, int(storedTxo[i].TxOutInd))
		if !res {
			println("Test failed: incorrect txo value")
			return
		}
	}

	for i := 0; i < len(storedUtxo); i++ {
		_, res := txo.GetUtxos(storedUtxo[i].OutTxHash, int(storedUtxo[i].TxOutInd))
		if !res {
			println("Test failed: incorrect utxo value")
			return
		}
	}

	println("TXO and UTXO values are correct")
	println("Test passed!")

}

func (rv *ReorganizationVerifTest) Launch() {
	rv.mcBlockAmmount = 5
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

	rv.genBlockTestAccounts(1)
	account.CurrAccount = account.Wallet[0]
	node_settings.Settings.RewardAddress = hashing.SHA1(account.CurrAccount.KeyPairList[0].PublKey)

	rv.genBlocks() // Generating blocks in mainchain
	time.Sleep(time.Millisecond * 200)

	rv.genAltchBlocks()
	rv.printResult()
}
