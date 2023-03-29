package networking

import (
	"bavovnacoin/node_controller/node_settings"
	"encoding/binary"
	"time"
)

func (l *Listener) SendNodeTime(reply *Reply) error {
	localTime := time.Now().Unix()
	var locTimeBytes []byte
	binary.LittleEndian.PutUint64(locTimeBytes, uint64(localTime))
	*reply = Reply{locTimeBytes}
	return nil
}

func (c *Connection) GetNodeTime() int64 {
	var repl Reply
	err := c.client.Call("Listener.SendNodeTime", nil, &repl)
	if err != nil {
		return -1
	}

	return int64(binary.BigEndian.Uint64(repl.Data))
}

func GetSettingsNodesTime() []int64 {
	var nodesTimeArray []int64
	var conn Connection
	var currAddrInd int = -1
	var allowConnect bool = true

	for allowConnect {
		allowConnect, currAddrInd = conn.EstablishAddresses(node_settings.Settings.OtherNodesAddresses, currAddrInd, "")
		nodeTime := conn.GetNodeTime()

		if nodeTime != -1 {
			nodesTimeArray = append(nodesTimeArray, nodeTime)
		}
	}
	return nodesTimeArray
}
