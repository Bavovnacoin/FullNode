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

func (c *Connection) EstablishAddresses(addresses []string, addrInd int) (bool, int) {
	for true {
		res := c.Establish(addresses[addrInd])
		if res {
			return true, addrInd
		} else if addrInd+1 < len(addresses[addrInd]) {
			addrInd++
			c.Close()
		} else {
			return false, addrInd
		}
	}
	return false, addrInd
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
