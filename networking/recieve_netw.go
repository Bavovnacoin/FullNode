package networking

import (
	"bavovnacoin/node/node_settings"
	"log"
	"net"
	"net/rpc"
)

var Inbound *net.TCPListener

type Listener int

type Reply struct {
	Data []byte
}

func StartRPCListener() (*net.TCPListener, bool, error) {
	addy, err := net.ResolveTCPAddr("tcp", node_settings.Settings.RPCip)
	if err != nil {
		log.Println(err)
		return nil, false, err
	}

	inbound, err := net.ListenTCP("tcp", addy)
	if err != nil {
		log.Println(err)
		return inbound, false, err
	}

	listener := new(Listener)
	rpc.Register(listener)
	go rpc.Accept(inbound)

	return inbound, true, err
}

func StopRPCListener() bool {
	err := Inbound.Close()
	if err == nil {
		return true
	}
	return false
}
