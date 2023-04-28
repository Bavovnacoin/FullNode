package networking_p2p

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/byteArr"
	"bavovnacoin/hashing"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
)

type BlockHashProposal struct {
	BlockHash byteArr.ByteArr
	Height    uint64
}

type BlockProposal struct {
	Block  blockchain.Block
	Height uint64
}

func TryAddCameBlock(data []byte, peerId peer.ID) bool {
	if data[0] == 3 { // Check hash for block existance
		var bhproposal BlockHashProposal
		byteArr.FromByteArr(data[1:], &bhproposal)

		height, _ := byteArr.ToByteArr(bhproposal.Height - 1)
		if !blockchain.IsBlockExists(bhproposal.BlockHash, bhproposal.Height) {
			return SendDataOnPeerId(append([]byte{4}, height...), peerId)
		}
		return false

	} else if data[0] == 4 {
		var height uint64
		byteArr.FromByteArr(data[1:], &height)
		block, isGotten := blockchain.GetBlock(height, 0)

		if !isGotten {
			return false
		}

		var blockToSend BlockProposal
		blockToSend.Height = height
		blockToSend.Block = block
		bytesToSend, _ := byteArr.ToByteArr(blockToSend)

		isSended := SendDataOnPeerId(append([]byte{5}, bytesToSend...), peerId)
		return isSended

	} else if data[0] == 5 {
		var newBlock BlockProposal
		isConv := byteArr.FromByteArr(data[1:], &newBlock)
		if !isConv {
			return false
		}

		nodesTime = []int64{}
		GetNodesTime()
		time.Sleep(50 * time.Millisecond)

		if blockchain.TryCameBlockToAdd(newBlock.Block, newBlock.Height, nodesTime) {
			ProposeNewBlock(newBlock.Block, newBlock.Height)
		}
	}
	return true
}

func ProposeNewBlock(block blockchain.Block, height uint64) bool {
	var blockHash byteArr.ByteArr
	blockHashString := hashing.SHA1(blockchain.BlockHeaderToString(block))
	blockHash.SetFromHexString(blockHashString, 20)

	var bhproposal BlockHashProposal
	bhproposal.BlockHash = blockHash
	bhproposal.Height = height

	hashProp, isConv := byteArr.ToByteArr(bhproposal)
	if !isConv {
		return false
	}

	return SendDataToAllConnectedPeers(append([]byte{3}, hashProp...))
}
