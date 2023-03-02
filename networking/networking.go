package networking

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

var myAddress string = "0.0.0.0:12345"

type Listener int

type Reply struct {
	Data []byte
}

func (l *Listener) PingPong(line []byte, reply *Reply) {
	rv := string(line)
	fmt.Printf("Receive: %v\n", rv)
	*reply = Reply{[]byte("pong")}
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
