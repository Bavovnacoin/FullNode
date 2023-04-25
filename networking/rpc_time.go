package networking

import (
	"bavovnacoin/node/node_controller/node_settings"
	"encoding/binary"
	"fmt"
	"time"
)

func (l *Listener) SendNodeTime(stubArg int, reply *Reply) error {
	localTime := time.Now().UTC().Unix()
	var locTimeBytes []byte = make([]byte, 8)
	binary.LittleEndian.PutUint64(locTimeBytes, uint64(localTime))
	*reply = Reply{locTimeBytes}
	return nil
}

func (c *Connection) GetNodeTime() int64 {
	var repl Reply
	err := c.client.Call("Listener.SendNodeTime", 0, &repl)
	if err != nil {
		fmt.Println(err)
		return -1
	}

	return int64(binary.LittleEndian.Uint64(repl.Data))
}

func GetSettingsNodesTime() []int64 {
	var nodesTimeArray []int64
	var conn Connection
	var currAddrInd int = -1
	var allowConnect bool = true

	for allowConnect {
		allowConnect, currAddrInd = conn.EstablishAddresses(node_settings.Settings.OtherNodesAddresses, currAddrInd, "")
		if !allowConnect {
			break
		}
		nodeTime := conn.GetNodeTime()

		if nodeTime != -1 {
			nodesTimeArray = append(nodesTimeArray, nodeTime)
		}
	}
	return nodesTimeArray
}
