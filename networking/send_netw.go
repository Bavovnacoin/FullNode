package networking

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/byteArr"
	"bavovnacoin/transaction"
	"bytes"
	"encoding/gob"
	"net/rpc"
)

type TxProposal struct {
	Tx      transaction.Transaction
	Address string
}

type BlockHashProposal struct {
	BlockHash byteArr.ByteArr
	Height    uint64
}

type BlockProposal struct {
	Block   blockchain.Block
	Address string
}

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

func (c *Connection) EstablishAddresses(addresses []string, addrInd int, avoidAddress string) (bool, int) {
	var res bool = false
	for addrInd < len(addresses) {
		if res {
			return true, addrInd
		} else if addrInd+1 < len(addresses) {
			addrInd++
			if addresses[addrInd] == avoidAddress {
				continue
			}
			res = c.Establish(addresses[addrInd])
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
