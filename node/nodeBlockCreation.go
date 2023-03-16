package node

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/byteArr"
	"bavovnacoin/node_controller/command_executor"
	"fmt"
	"log"
)

var AllowCreateBlock bool = true
var CreatedBlock blockchain.Block

func AddBlock(allowLogPrint bool) bool {
	if AllowCreateBlock {
		if allowLogPrint {
			log.Println("Creating a new block")
		}
		go CreateBlockLog(blockchain.GetBits(allowLogPrint), allowLogPrint)
		AllowCreateBlock = false
	}

	if CreatedBlock.MerkleRoot != "" { // Is block mined check
		isBlockValid := blockchain.ValidateBlock(CreatedBlock, int(blockchain.BcLength), true, false)
		AddBlockLog(allowLogPrint, isBlockValid)
		CreatedBlock.MerkleRoot = ""
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
	newBlock = blockchain.MineBlock(newBlock, 1, allowPrint)
	CreatedBlock = newBlock
	command_executor.PauseCommand()
}

func AddBlockLog(allowPrint bool, isBlockValid bool) bool {
	isBlockAdded := false

	if isBlockValid {
		blockchain.AddBlockToBlockchain(CreatedBlock, true)
		if allowPrint {
			log.Println("Block is added to blockchain. Current height: " + fmt.Sprint(blockchain.BcLength+1) + "\n")
		}
		blockchain.IncrBcHeight()
		isBlockAdded = true
	} else {
		if allowPrint {
			log.Println("Block is not added\n")
		}
		isBlockAdded = false
	}

	AllowCreateBlock = true
	CreatedBlock.MerkleRoot = ""
	command_executor.PauseCommand()
	return isBlockAdded
}
