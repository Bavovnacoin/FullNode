package node

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/byteArr"
	"bavovnacoin/networking"
	"bavovnacoin/node_controller/command_executor"
	"fmt"
	"log"
	"time"
)

func AddBlock(allowLogPrint bool) bool {
	for !blockchain.IsMiningDone { //Waiting for mining to be done
		time.Sleep(1 * time.Millisecond)
	}

	if blockchain.AllowCreateBlock {
		if allowLogPrint {
			log.Println("Creating a new block")
		}
		go CreateBlockLog(blockchain.GetBits(allowLogPrint), allowLogPrint)
		blockchain.AllowCreateBlock = false
	}

	if blockchain.CreatedBlock.MerkleRoot != "" { // Is block mined check
		isBlockValid := blockchain.VerifyBlock(blockchain.CreatedBlock, int(blockchain.BcLength), true, false)
		AddBlockLog(allowLogPrint, isBlockValid)
		blockchain.CreatedBlock.MerkleRoot = ""
		return true
	}
	command_executor.PauseCommand()
	return false
}

func CreateBlockLog(bits uint64, allowPrint bool) {
	var rewardAdr byteArr.ByteArr
	rewardAdr.SetFromHexString(blockchain.RewardAddress, 20)
	newBlock := blockchain.CreateBlock(rewardAdr, allowPrint)
	newBlock.Bits = bits
	newBlock.Chainwork = blockchain.GetChainwork(newBlock)
	var miningRes bool
	newBlock, miningRes = blockchain.MineBlock(newBlock, 1, allowPrint)

	if !miningRes {
		blockchain.AllowCreateBlock = true
		return
	}

	blockchain.IsMiningDone = true
	blockchain.RemoveTxsFromMempool(newBlock.Transactions[1:])
	blockchain.CreatedBlock = newBlock
	command_executor.PauseCommand()
}

func AddBlockLog(allowPrint bool, isBlockValid bool) bool {
	isBlockAdded := false

	if isBlockValid {
		blockAddRes := blockchain.AddBlockToBlockchain(blockchain.CreatedBlock, 0)
		if !blockAddRes {
			return false
		}

		if allowPrint {
			log.Println("Block is added to blockchain. Current height: " + fmt.Sprint(blockchain.BcLength+1) + "\n")
		}
		blockchain.IncrBcHeight()
		isBlockAdded = true
		println()
		networking.ProposeBlockToSettingsNodes(blockchain.CreatedBlock, "")
	} else {
		if allowPrint {
			log.Println("Created block is not valid.")
			println()
		}
		isBlockAdded = false
	}

	blockchain.AllowCreateBlock = true
	blockchain.CreatedBlock.MerkleRoot = ""
	command_executor.PauseCommand()
	return isBlockAdded
}
