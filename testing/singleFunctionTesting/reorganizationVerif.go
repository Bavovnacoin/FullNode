/*
	Checks outputs and blockchain structure after reorganization
*/

package singleFunctionTesting

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/hashing"
	"bavovnacoin/networking_p2p"
	"bavovnacoin/node/node_controller/command_executor"
	"bavovnacoin/node/node_controller/node_settings"
	"bavovnacoin/node/node_validator"
	"bavovnacoin/testing/account"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type ReorganizationVerifTest struct {
	SingleFunctionTesting
	mcBlockAmmount uint64
	acBlockAmmount uint64

	source rand.Source
	random *rand.Rand
}

func (rv *ReorganizationVerifTest) genBlocks() {
	command_executor.ComContr.FullNodeWorking = true
	go rv.nodeWorkListener(rv.mcBlockAmmount)
	go rv.txForming("abc", rv.random)
	node_validator.BlockGen(false)
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
		blockchain.TryCameBlockToAdd(blockchain.CreatedBlock, blockchain.BcLength-1+uint64(i), otherNodesTime, true)
		prevHash = hashing.SHA1(blockchain.BlockHeaderToString(blockchain.CreatedBlock))
		bl = blockchain.CreatedBlock
	}
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

	result := rv.checkMcTxo()
	result.PrintTestResult()
}

func (rv *ReorganizationVerifTest) Launch() {
	rv.mcBlockAmmount = 5
	rv.acBlockAmmount = 3
	node_settings.Settings.GetSettings()
	networking_p2p.StartP2PCommunication()

	InitTestDb(true)

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
