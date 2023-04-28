package networking

import (
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
