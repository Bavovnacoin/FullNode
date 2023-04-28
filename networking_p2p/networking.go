package networking_p2p

import (
	"bavovnacoin/node/node_controller/node_settings"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/multiformats/go-multiaddr"
)

var Peer host.Host
var OtherPeersIds []peer.ID

const PROTOCOL_ID = protocol.ID("/bvc/1.0.0")

func StreamHandler(s network.Stream) {
	data, err := ioutil.ReadAll(s)
	if err != nil {
		log.Fatal(err)
	}

	peerID := s.Conn().RemotePeer()

	TryHandleSynchronization(data, peerID)
}

// My address (localhost) /ip4/127.0.0.1/tcp/58818
func StartP2PCommunication() {
	peerIdInd = 0
	pkBytes := node_settings.Settings.GetPrivKey()
	privKey, _ := crypto.UnmarshalSecp256k1PrivateKey(pkBytes)

	var err error
	Peer, err = libp2p.New(
		libp2p.Identity(privKey),
		libp2p.ListenAddrStrings(node_settings.Settings.MyAddress),
	)

	if err == nil {
		Peer.SetStreamHandler(PROTOCOL_ID, StreamHandler)
		fmt.Printf("Started peer on address %s/%s\n", Peer.Addrs()[0], Peer.ID().Pretty())
	} else {
		fmt.Println("Unable to start a peer.", err)
	}

	addSettingsAddresses()
	//SendOutNewBlock()
}

func addOtherAddress(address string) bool {
	arr := strings.Split(address, "/")

	maddr, err := multiaddr.NewMultiaddr(strings.Join(arr[:len(arr)-1], "/"))
	if err != nil {
		return false
	}

	id, res := peer.Decode(arr[len(arr)-1])
	if res != nil {
		return false
	}

	Peer.Peerstore().AddAddrs(id, []multiaddr.Multiaddr{maddr}, peerstore.PermanentAddrTTL)
	OtherPeersIds = append(OtherPeersIds, id)
	return true
}

func addSettingsAddresses() {
	for i := 0; i < len(node_settings.Settings.OtherNodesAddresses); i++ {
		addOtherAddress(node_settings.Settings.OtherNodesAddresses[i])
	}
}

func SendDataToAll(data []byte) bool {
	for _, id := range OtherPeersIds {
		SendDataOnPeerId(data, id)
	}
	return true
}

func SendDataOnPeerId(data []byte, id peer.ID) bool {
	if err := Peer.Connect(context.Background(), Peer.Peerstore().PeerInfo(id)); err == nil {
		stream, err := Peer.NewStream(context.Background(), id, PROTOCOL_ID)
		if err != nil {
			return false
		}

		if _, err := stream.Write(data); err != nil {
			return false
		}

		if err := stream.Close(); err != nil {
			return false
		}
		return true
	}
	return false
}
