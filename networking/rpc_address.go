package networking

import (
	"bavovnacoin/byteArr"
	"bavovnacoin/txo"
)

func (l *Listener) IsAddrExist(address []byte, reply *Reply) error {
	var addr byteArr.ByteArr = byteArr.ByteArr{ByteArr: address}
	isAddrExist := txo.IsOutAddrExist(addr)

	if isAddrExist {
		*reply = Reply{[]byte{1}}
	} else {
		*reply = Reply{[]byte{0}}
	}
	return nil
}
