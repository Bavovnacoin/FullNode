package networking_p2p

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/byteArr"

	"github.com/libp2p/go-libp2p/core/peer"
)

var peerIdInd int

func incrPeerId() int {
	return (peerIdInd + 1) % len(OtherPeersIds)
}

func TryHandleSynchronization(data []byte, peerId peer.ID) bool {
	if data[0] == 1 { // Request handler
		var height uint64
		isConv := byteArr.FromByteArr(data[1:], &height)
		if !isConv {
			return false
		}

		blocks := GetBlocksOnHeight(height)

		reqData, isReqConv := byteArr.ToByteArr(blocks)
		if !isReqConv {
			return false
		}
		data = append([]byte{2}, reqData...)
		SendDataOnPeerId(data, peerId)
	} else if data[0] == 2 { // Responce handler
		var request BlockRequest
		isConv := byteArr.FromByteArr(data[1:], &request)
		if !isConv {
			return false
		}

		isCurrSyncSuccess := SyncAddBlocks(request.Blocks)
		if !isCurrSyncSuccess {
			if peerIdInd+1 < len(OtherPeersIds) {
				peerIdInd++
			} else {
				IsSyncEnded = true
				IsSyncSuccess = false
			}
		}

		if request.IsMoreBlocks {
			RequestBlocks(blockchain.BcLength, OtherPeersIds[peerIdInd])
		} else {
			IsSyncEnded = true
			IsSyncSuccess = true
			return true
		}
	}
	return true
}
