package networking

import (
	"bytes"
	"encoding/gob"
	"net/rpc"
)

type Connection struct {
	client *rpc.Client
}

func (c *Connection) Establish(address string) bool {
	var err error
	c.client, err = rpc.Dial("tcp", address)
	if err != nil {
		return false
	}
	return true
}

func (c *Connection) Close() {
	c.client.Close()
}

func (c *Connection) ToByteArr(data any) ([]byte, bool) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(data)
	if err != nil {
		return nil, false
	}

	return buffer.Bytes(), true
}
