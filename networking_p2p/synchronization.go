package networking_p2p

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/byteArr"
	"bavovnacoin/hashing"
	"bavovnacoin/synchronization"
	"log"

	"github.com/libp2p/go-libp2p/core/peer"
)

type BlocksOnHeight struct {
	Blocks []blockchain.BlockChainId
	Height uint64
}

type BlockRequest struct {
	BcHeight     uint64
	Blocks       []BlocksOnHeight
	IsMoreBlocks bool
}

const retrStep uint64 = 5

var IsSyncEnded bool = false
var IsSyncSuccess bool = false

func (pd *PeerData) RequestBlocks(startFromHeight uint64, peerId peer.ID) bool {
	startByteArr, isConv := byteArr.ToByteArr(startFromHeight)
	if !isConv {
		return false
	}

	data := []byte{1}
	data = append(data, startByteArr...)
	isSended := pd.SendDataOnPeerId(data, peerId)
	return isSended
}

func GetBlocksOnHeight(startFromBlock uint64) BlockRequest {
	var request BlockRequest
	i := startFromBlock
	for ; i < startFromBlock+retrStep && i < blockchain.BcLength; i++ {
		blocks, res := blockchain.GetBlocksOnHeight(i)
		if res {
			request.Blocks = append(request.Blocks, BlocksOnHeight{Blocks: blocks, Height: i})
		} else {
			break
		}
	}

	request.BcHeight = blockchain.BcLength
	if i < blockchain.BcLength {
		request.IsMoreBlocks = true
	} else {
		request.IsMoreBlocks = false
	}

	return request
}

func SyncAddBlocks(blocks []BlocksOnHeight, allowLogging bool) bool {
	checkpCorresp := false
	addCount := 0
	i := 0
	for ; i < len(blocks); i++ {
		for _, bl := range blocks[i].Blocks {
			checkpCorresp = synchronization.CheckForBlockCorrespondence(blocks[i].Height, bl.Block)
			if checkpCorresp {
				bcBlock, isBlockGotten := blockchain.GetBlock(blocks[i].Height, bl.ChainId)
				if !isBlockGotten ||
					hashing.SHA1(blockchain.BlockHeaderToString(bcBlock)) != hashing.SHA1(blockchain.BlockHeaderToString(bl.Block)) {
					blockchain.AddBlockToBlockchain(bl.Block, blocks[i].Height, bl.ChainId, bl.ChainId == 0)
					blockchain.IncrBcHeight(bl.ChainId)
					if bl.ChainId == 0 {
						blockchain.LastBlock = bl.Block
						addCount++
					}
				}
			} else {
				if allowLogging {
					log.Println("Address sent an incorrect block")
				}
				return false
			}
		}
	}

	if allowLogging {
		log.Printf("Checked %d levels. Current bc height: %d\n", len(blocks), blockchain.BcLength)
	}
	return true
}

func (pd *PeerData) StartSync() bool {
	if len(pd.OtherPeersIds) > 0 {
		return pd.RequestBlocks(0, pd.OtherPeersIds[0])
	}
	return false
}
