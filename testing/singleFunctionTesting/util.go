package singleFunctionTesting

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/byteArr"
	"bavovnacoin/dbController"
	"bavovnacoin/hashing"
	"bavovnacoin/node/node_controller/command_executor"
	"bavovnacoin/node/node_settings"
	"bavovnacoin/testing/account"
	"bavovnacoin/transaction"
	"bavovnacoin/txo"
	"fmt"
	"math/rand"
	"os"
)

type SingleFunctionTesting struct{}

func (sft *SingleFunctionTesting) genBlockTestAccounts(ammount int) {
	for i := 0; i < ammount; i++ {
		account.Wallet = append(account.Wallet, account.GenAccount(fmt.Sprint(len(account.Wallet))))
	}
}

func (sft *SingleFunctionTesting) CreateTestTx(textForAddress string, random *rand.Rand) (transaction.Transaction, bool) {
	fee := random.Intn(5) + 1
	isGenLocktime := random.Intn(5)
	var locktime uint
	if isGenLocktime == 2 {
		locktime = uint(int(blockchain.BcLength+1) + random.Intn(2) + 1)
	}

	var outAddr []byteArr.ByteArr
	outAddr = append(outAddr, byteArr.ByteArr{})
	outAddr[0].SetFromHexString(hashing.SHA1(textForAddress), 20)

	var outValue []uint64
	outValue = append(outValue, 1000)

	tx, isValid := transaction.CreateTransaction(fmt.Sprint(0), outAddr, outValue, fee, locktime)

	return tx, isValid == ""
}

// Generates transactions and adds them to the mempool
func (sft *SingleFunctionTesting) txForming(address string, random *rand.Rand) {
	for command_executor.ComContr.FullNodeWorking {
		tx, isValid := sft.CreateTestTx(address, random)
		if !isValid {
			continue
		}
		blockchain.AddTxToMempool(tx, true)
	}
}

func (sft *SingleFunctionTesting) nodeWorkListener(blocksCount uint64) {
	for true {
		if blockchain.BcLength >= blocksCount {
			command_executor.ComContr.FullNodeWorking = false
			return
		}
	}
}

func (sft *SingleFunctionTesting) getOutputsFromMainchain() ([]txo.TXO, []txo.TXO, bool) {
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

// make a method that prints result and a method that prints test output

type TxoTestCheckRes struct {
	// Values from the database
	storedTxo  []txo.TXO
	storedUtxo []txo.TXO

	// Values from the mainchain
	mcTxo  []txo.TXO
	mcUtxo []txo.TXO

	value int // Test result
}

func (ttch *TxoTestCheckRes) PrintTestResult() {
	switch ttch.value {
	case 0:
		println("TXO and UTXO values are correct")
		println("Test passed!")
	case 1:
		println("Error when getting values from a mainchain")
	case 2:
		println("Test failed. Wrong txo ammount.")
	case 3:
		println("Test failed. Wrong гtxo ammount.")
	case 4:
		println("Test failed: incorrect txo value")
	case 5:
		println("Test failed: incorrect utxo value")
	}
}

func (ttch *TxoTestCheckRes) PrintTestOutput() {
	fmt.Printf("TXO: stored: %d, mc: %d\n", len(ttch.storedTxo), len(ttch.mcTxo))
	fmt.Printf("UTXO: stored: %d, mc: %d\n", len(ttch.storedUtxo), len(ttch.mcUtxo))
}

func (sft *SingleFunctionTesting) checkMcTxo() TxoTestCheckRes {
	testRes := TxoTestCheckRes{}
	testRes.value = 0

	mTxo, mUtxo, res := sft.getOutputsFromMainchain()
	storedTxo, _ := txo.GetTxoList("txo")
	storedUtxo, _ := txo.GetTxoList("utxo")

	testRes.storedTxo = storedTxo
	testRes.storedUtxo = storedUtxo
	testRes.mcTxo = mTxo
	testRes.mcUtxo = mUtxo

	if !res {
		testRes.value = 1
		return testRes
	}

	if len(mTxo) != len(storedTxo) {
		println(len(mTxo), len(storedTxo))
		println("Test failed. Wrong txo ammount.")
		testRes.value = 2
		return testRes
	} else if len(mUtxo) != len(storedUtxo) {
		println("Test failed. Wrong utxo ammount.")
		testRes.value = 3
		return testRes
	}

	for i := 0; i < len(storedTxo); i++ {
		_, res := txo.GetTxos(storedTxo[i].OutTxHash, int(storedTxo[i].TxOutInd))
		if !res {
			println("Test failed: incorrect txo value")
			testRes.value = 4
			return testRes
		}
	}

	for i := 0; i < len(storedUtxo); i++ {
		_, res := txo.GetUtxos(storedUtxo[i].OutTxHash, int(storedUtxo[i].TxOutInd))
		if !res {
			println("Test failed: incorrect utxo value")
			testRes.value = 5
			return testRes
		}
	}
	return testRes
}

func (sft *SingleFunctionTesting) CreateBlock(bits uint64, prevHash string, lastBlock blockchain.Block, allowPrint bool) (blockchain.Block, bool) {
	var rewardAdr byteArr.ByteArr
	rewardAdr.SetFromHexString(node_settings.Settings.RewardAddress, 20)
	newBlock := blockchain.CreateBlock(rewardAdr, prevHash, allowPrint)
	newBlock.Bits = bits
	newBlock.Chainwork = blockchain.GetChainwork(newBlock, lastBlock)
	var miningRes bool

	newBlock, miningRes = blockchain.MineBlock(newBlock, 1, allowPrint)

	if !miningRes {
		blockchain.AllowCreateBlock = true
		return newBlock, false
	}

	blockchain.IsMiningDone = true
	blockchain.RemoveTxsFromMempool(newBlock.Transactions[1:])
	blockchain.CreatedBlock = newBlock
	return newBlock, true
}

func InitTestDb(allowLogging bool) {
	dbController.DbPath = "testing/testData"
	if _, err := os.Stat(dbController.DbPath); err == nil {
		os.RemoveAll(dbController.DbPath)

		if allowLogging {
			println("Removed test db from a previous test.")
		}
	}
	dbController.DB.OpenDb()
}
