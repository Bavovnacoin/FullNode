package networking

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/byteArr"
	"bavovnacoin/hashing"
	"bavovnacoin/node/node_controller/node_settings"
	"log"
)

const retrStep uint64 = 4

type BlocksOnHeight struct {
	Blocks []blockchain.BlockChainId
	Height uint64
}

type BlockRequest struct {
	BcHeight uint64
	Blocks   []BlocksOnHeight
}

// TODO: chaneg according to the altchain idea
func (l *Listener) SendBlocks(startFromBlock []byte, reply *Reply) error {
	var startBlockToRetrieve uint64
	byteArr.FromByteArr(startFromBlock, &startBlockToRetrieve)
	var request BlockRequest

	for i := startBlockToRetrieve; i < startBlockToRetrieve+retrStep && i < blockchain.BcLength; i++ {
		blocks, res := blockchain.GetBlocksOnHeight(i)
		if res {
			request.Blocks = append(request.Blocks, BlocksOnHeight{Blocks: blocks, Height: i})
		} else {
			break
		}
	}
	request.BcHeight = blockchain.BcLength

	requestStructByte, _ := byteArr.ToByteArr(request)
	*reply = Reply{requestStructByte}
	return nil
}

func (c *Connection) RequestBlocks(startFromHeight uint64) ([]BlocksOnHeight, uint64, bool) {
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

func (l *Listener) AddProposedBlock(blockProposalByteArr []byte, reply *Reply) error {
	var blockProp BlockProposal
	byteArr.FromByteArr(blockProposalByteArr, &blockProp)

	if blockchain.TryCameBlockToAdd(blockProp.Block, blockProp.Height, GetSettingsNodesTime()) {
		*reply = Reply{[]byte{1}}
		ProposeBlockToSettingsNodes(blockProp.Block, blockProp.Address)
	} else {
		*reply = Reply{[]byte{0}}
	}
	return nil
}

func (l *Listener) GetBlockProposal(blockHashPropByteArr []byte, reply *Reply) error {
	var blockHashProposal BlockHashProposal
	byteArr.FromByteArr(blockHashPropByteArr, &blockHashProposal)

	if !blockchain.IsBlockExists(blockHashProposal.BlockHash, blockHashProposal.Height) {
		*reply = Reply{[]byte{1}}
	} else {
		*reply = Reply{[]byte{0}}
		return nil
	}
	return nil
}

// TODO: make in a single function (probability of using second funxtion directly)
func (c *Connection) ProposeBlockToOtherNode(blockHash []byte, block blockchain.Block, blockHeight uint64) bool {
	var repl Reply
	var blockHashProposal BlockHashProposal
	blockHashProposal.BlockHash.ByteArr = blockHash
	blockHashProposal.Height = blockHeight
	bhPropBytes, _ := c.ToByteArr(blockHashProposal)

	err := c.client.Call("Listener.GetBlockProposal", bhPropBytes, &repl)
	if err != nil {
		return false // Problem when accessing an RPC function
	}
	if repl.Data[0] == 1 {
		repl.Data = []byte{}

		var blockProp BlockProposal
		blockProp.Block = block
		blockProp.Address = node_settings.Settings.MyAddress
		blockProp.Height = blockHeight
		propBytes, _ := c.ToByteArr(blockProp)

		println("Sent block with hash", hashing.SHA1(blockchain.BlockHeaderToString(block)))

		err := c.client.Call("Listener.AddProposedBlock", propBytes, &repl)
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

		connection.ProposeBlockToOtherNode(blockHash.ByteArr, block, blockchain.BcLength-1)
	}
	return true
}
