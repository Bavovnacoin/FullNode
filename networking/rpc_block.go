package networking

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/byteArr"
	"bavovnacoin/hashing"
	"bavovnacoin/node_controller/node_settings"
	"log"
)

const retrStep uint64 = 4

type BlockRequest struct {
	BcHeight uint64
	Blocks   []blockchain.Block
}

func (l *Listener) SendBlocks(startFromBlock []byte, reply *Reply) error {
	var startBlockToRetrieve uint64
	byteArr.FromByteArr(startFromBlock, &startBlockToRetrieve)
	var request BlockRequest

	for i := startBlockToRetrieve; i < startBlockToRetrieve+retrStep && i < blockchain.BcLength; i++ {
		block, res := blockchain.GetBlock(i)
		if res {
			request.Blocks = append(request.Blocks, block)
		} else {
			break
		}
	}
	request.BcHeight = blockchain.BcLength

	requestStructByte, _ := byteArr.ToByteArr(request)
	*reply = Reply{requestStructByte}
	return nil
}

// TODO: add a parameter of ammount of blocks to request (where 0 - max)
func (c *Connection) RequestBlocks(startFromHeight uint64) ([]blockchain.Block, uint64, bool) {
	var blockReq BlockRequest
	startByteArr, isConv := c.ToByteArr(startFromHeight)
	if !isConv {
		return blockReq.Blocks, blockReq.BcHeight, false
	}

	var repl Reply
	err := c.client.Call("Listener.SendBlocks", startByteArr, &repl)
	if err != nil {
		log.Println(err)
		return blockReq.Blocks, blockReq.BcHeight, false
	}

	byteArr.FromByteArr(repl.Data, &blockReq)

	return blockReq.Blocks, blockReq.BcHeight, true
}

func (l *Listener) AddProposedBlockToMemp(blockProposalByteArr []byte, reply *Reply) error {
	var blockProp BlockProposal
	byteArr.FromByteArr(blockProposalByteArr, &blockProp)

	if blockchain.TryCameBlockToAdd(blockProp.Block, GetSettingsNodesTime()) {
		*reply = Reply{[]byte{1}}
	} else {
		*reply = Reply{[]byte{0}}
	}
	return nil
}

func (l *Listener) GetBlockProposal(blockHashByteArr []byte, reply *Reply) error {
	var blockHash byteArr.ByteArr
	blockHash.ByteArr = blockHashByteArr
	var lastBlockHash byteArr.ByteArr
	lastBlockHash.SetFromHexString(hashing.SHA1(blockchain.BlockHeaderToString(blockchain.LastBlock)), 20)
	if !lastBlockHash.IsEqual(blockHash) {
		*reply = Reply{[]byte{1}}
	} else {
		*reply = Reply{[]byte{0}}
		return nil
	}
	return nil
}

func (c *Connection) ProposeBlockToOtherNode(blockHash []byte, block blockchain.Block) bool {
	var repl Reply
	err := c.client.Call("Listener.GetBlockProposal", blockHash, &repl)
	if err != nil {
		return false // Problem when accessing an RPC function
	}
	if repl.Data[0] == 1 {
		repl.Data = []byte{}

		var blockProp BlockProposal
		blockProp.Block = block
		blockProp.Address = node_settings.Settings.MyAddress
		propBytes, _ := c.ToByteArr(blockProp)

		err := c.client.Call("Listener.AddProposedBlockToMemp", propBytes, &repl)
		if err != nil || repl.Data[0] == 0 {
			return false // The node reverted this block
		}
	} else {
		return false // The node is already has this block
	}

	return true // No problems
}

func ProposeBlockToSettingsNodes(block blockchain.Block, avoidAddress string) bool {
	var blockHash byteArr.ByteArr
	blockHashString := hashing.SHA1(blockchain.BlockHeaderToString(block))
	blockHash.SetFromHexString(blockHashString, 20)

	var connection Connection
	var isNodesAccessible bool

	for i := 0; i < len(node_settings.Settings.OtherNodesAddresses); i++ {
		isNodesAccessible, i = connection.EstablishAddresses(node_settings.Settings.OtherNodesAddresses, i-1, avoidAddress)

		if !isNodesAccessible {
			return false
		}

		connection.ProposeBlockToOtherNode(blockHash.ByteArr, block)
	}
	return true
}
