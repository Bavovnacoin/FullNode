package dbController

import (
	"bavovnacoin/blockchain"
	"bytes"
	"encoding/gob"
)

func ToByteArr(data any) ([]byte, bool) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(data)
	if err != nil {
		return nil, false
	}

	return buffer.Bytes(), true
}

func FromByteArr(data []byte) any {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	var decData blockchain.Block
	err := dec.Decode(&decData)
	if err != nil {
		return decData
	}
	return decData
}
