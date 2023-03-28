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

func ProposeBlockToSettingsNodes(block blockchain.Block, avoidAddress string) bool {
	var blockHash byteArr.ByteArr
	txHashString := hashing.SHA1(blockchain.BlockHeaderToString(block))
	txHash.SetFromHexString(txHashString, 20)

	var connection Connection
	var isNodesAccessible bool

	for i := 0; i < len(node_settings.Settings.OtherNodesAddresses); i++ {
		isNodesAccessible, i = connection.EstablishAddresses(node_settings.Settings.OtherNodesAddresses, i-1, avoidAddress)

		if !isNodesAccessible {
			return false
		}

		connection.ProposeTxToOtherNode(txHash.ByteArr, tx)
	}
	return true
}
