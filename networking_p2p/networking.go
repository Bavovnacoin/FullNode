package networking_p2p

import (
	"bavovnacoin/node/node_controller/node_settings"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
)

var Peer host.Host

const PROTOCOL_ID = protocol.ID("/bvc/1.0.0")

func StreamHandler(s network.Stream) {
	log.Println("Got a new stream!")
}

// My address (localhost) /ip4/127.0.0.1/tcp/58818
func StartP2PCommunication() {
	pkBytes := node_settings.Settings.GetPrivKey()
	privKey, _ := crypto.UnmarshalSecp256k1PrivateKey(pkBytes)

	Peer, err := libp2p.New(
		libp2p.Identity(privKey),
		libp2p.ListenAddrStrings(node_settings.Settings.MyAddress),
	)

	if err == nil {
		Peer.SetStreamHandler(PROTOCOL_ID, StreamHandler)
		fmt.Printf("Started peer on address %s/%s\n", Peer.Addrs()[0], Peer.ID().Pretty())
	} else {
		fmt.Println("Unable to start a peer.", err)
	}
}
