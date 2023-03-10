package loadtesting

import (
	"bavovnacoin/transaction"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net/rpc"
)

type Connection struct {
	client *rpc.Client
}

type Reply struct {
	Data []byte
}

func (c *Connection) Establish() bool {
	var err error
	c.client, err = rpc.Dial("tcp", "localhost:12345")
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

func (c *Connection) FromByteArr(dataByte []byte, data interface{}) bool {
	buf := bytes.NewBuffer(dataByte)
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(data)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func (c *Connection) SendTransaction(tx transaction.Transaction, isAccepted *bool) bool {
	byteArr, isConv := c.ToByteArr(tx)
	if !isConv {
		return false
	}

	var repl Reply
	err := c.client.Call("Listener.AddNewTxToMemp", byteArr, &repl)
	if err != nil {
		log.Println(err)
		return false
	}

	if repl.Data[0] == 1 {
		*isAccepted = true
	} else {
		*isAccepted = false
	}
	return true
}
