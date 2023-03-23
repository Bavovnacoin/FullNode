package networking

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/byteArr"
	"log"
)

const retrStep uint64 = 10

func (l *Listener) GetBlocks(startFromBlock []byte, reply *Reply) error {
	var startBlockToRetrieve uint64
	byteArr.FromByteArr(startFromBlock, &startBlockToRetrieve)

	var blocksToSend []blockchain.Block
	for i := startBlockToRetrieve; i < startBlockToRetrieve+retrStep && i < blockchain.BcLength; i++ {
		block, res := blockchain.GetBlock(i)
		if res {
			blocksToSend = append(blocksToSend, block)
		} else {
			break
		}
	}

	blocksToSendByte, _ := byteArr.ToByteArr(blocksToSend)
	*reply = Reply{blocksToSendByte}
	return nil
}

func (c *Connection) RequestBlocks(startFromHeight uint64) ([]blockchain.Block, bool) {
	var blocks []blockchain.Block
	startByteArr, isConv := c.ToByteArr(startFromHeight)
	if !isConv {
		return blocks, false
	}

	var repl Reply
	err := c.client.Call("Listener.GetBlocks", startByteArr, &repl)
	if err != nil {
		log.Println(err)
		return blocks, false
	}

	byteArr.FromByteArr(repl.Data, &blocks)
	return blocks, true
}
