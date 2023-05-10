package networking_p2p

import (
	"bavovnacoin/byteArr"
	"encoding/binary"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
)

var nodesTime []int64

func (pd *PeerData) TryHandleTime(data []byte, peerId peer.ID) bool {
	if data[0] == 6 {
		locTime := time.Now().UTC().Unix()
		timeByteArr, _ := byteArr.ToByteArr(locTime)
		return pd.SendDataOnPeerId(append([]byte{7}, timeByteArr...), peerId)
	} else if data[0] == 7 {
		var nodeTime int64
		byteArr.FromByteArr(data[1:], &nodeTime)
		nodesTime = append(nodesTime, nodeTime)
		return true
	}
	return false
}

func (pd *PeerData) GetNodesTime() bool {
	localTime := time.Now().UTC().Unix()
	var locTimeBytes []byte = make([]byte, 8)
	binary.LittleEndian.PutUint64(locTimeBytes, uint64(localTime))
	return pd.SendDataToAllConnectedPeers(append([]byte{7}, locTimeBytes...))
}
