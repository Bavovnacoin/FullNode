package dbController

import (
	"bytes"
	"encoding/gob"
	"fmt"
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

func FromByteArr(dataByte []byte, data interface{}) bool {
	buf := bytes.NewBuffer(dataByte)
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(data)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
