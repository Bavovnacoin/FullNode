package loadtesting

import (
	"bavovnacoin/byteArr"
	"bavovnacoin/networking"
	"bavovnacoin/node_controller/node_settings"
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
	node_settings.Settings.GetSettings()
	c.client, err = rpc.Dial("tcp", "localhost:8080") //node_settings.Settings.MyAddress
	if err != nil {
		log.Println(err)
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

	var l networking.Listener
	var repl networking.Reply
	err := l.AddNewTxToMemp(byteArr, &repl)
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

func (c *Connection) IsAddrExist(addr byteArr.ByteArr) bool {
	var repl Reply
	err := c.client.Call("Listener.IsAddrExist", addr.ByteArr, &repl)
	if err != nil {
		log.Println(err)
		return false
	}

	if repl.Data[0] == 0 {
		return false
	}

	return true
}

// GetMyUtxo - in wallet
func (c *Connection) GetUtxoByAddress(addresses []byteArr.ByteArr) bool {
	byteArr, isConv := c.ToByteArr(addresses)
	if !isConv {
		return false
	}

	var repl networking.Reply
	var l networking.Listener
	err := l.GetUtxoByAddr(byteArr, &repl)
	if err != nil {
		log.Println(err)
		return false
	}

	return true

}
