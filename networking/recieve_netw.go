package networking

import (
	"log"
	"net"
	"net/rpc"
)

var myAddress string = "localhost:12345"

type Listener int

type Reply struct {
	Data []byte
}

func StartRPCListener() (bool, error) {
	addy, err := net.ResolveTCPAddr("tcp", myAddress)
	if err != nil {
		log.Println(err)
		return false, err
	}

	inbound, err := net.ListenTCP("tcp", addy)
	if err != nil {
		log.Println(err)
		return false, err
	}

	listener := new(Listener)
	rpc.Register(listener)
	go rpc.Accept(inbound)

	return true, err
}