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

var Peer PeerData

type PeerData struct {
	Peer          host.Host
	OtherPeersIds []peer.ID
}

const PROT_NAME = "/bvc/1.0.0"
const PROTOCOL_ID = protocol.ID(PROT_NAME)

func (pd *PeerData) StreamHandler(s network.Stream) {
	data, err := ioutil.ReadAll(s)
	if err != nil {
		log.Fatal(err)
	}

	peerID := s.Conn().RemotePeer()

	pd.TryHandleSynchronization(data, peerID)
	pd.TryAddCameBlock(data, peerID)
	pd.TryHandleTime(data, peerID)
	pd.TryHandleTx(data, peerID)
}

// My address (localhost) /ip4/127.0.0.1/tcp/58818
func (pd *PeerData) StartP2PCommunication() {
	peerIdInd = 0
	pkBytes := node_settings.Settings.GetPrivKey()
	privKey, _ := crypto.UnmarshalSecp256k1PrivateKey(pkBytes)

	var err error
	pd.Peer, err = libp2p.New(
		libp2p.Identity(privKey),
		libp2p.ListenAddrStrings(node_settings.Settings.MyAddress),
	)

	if err == nil {
		pd.Peer.SetStreamHandler(PROTOCOL_ID, pd.StreamHandler)
		fmt.Printf("Started peer on address %s/%s\n", pd.Peer.Addrs()[0], pd.Peer.ID().Pretty())
	} else {
		fmt.Println("Unable to start a peer.", err)
	}

	pd.addSettingsAddresses()
}

func (pd *PeerData) addOtherAddress(address string) bool {
	arr := strings.Split(address, "/")

	maddr, err := multiaddr.NewMultiaddr(strings.Join(arr[:len(arr)-1], "/"))
	if err != nil {
		return false
	}

	id, res := peer.Decode(arr[len(arr)-1])
	if res != nil {
		return false
	}

	pd.Peer.Peerstore().AddAddrs(id, []multiaddr.Multiaddr{maddr}, peerstore.PermanentAddrTTL)
	pd.OtherPeersIds = append(pd.OtherPeersIds, id)
	return true
}

func (pd *PeerData) addSettingsAddresses() {
	for i := 0; i < len(node_settings.Settings.OtherNodesAddresses); i++ {
		pd.addOtherAddress(node_settings.Settings.OtherNodesAddresses[i][0])
	}
}

func (pd *PeerData) SendDataToAllConnectedPeers(data []byte) bool {
	activePeersCounter := 0
	peerIds := pd.Peer.Peerstore().Peers()

	for i := 0; i < len(peerIds); i++ {
		if pd.SendDataOnPeerId(data, peerIds[i]) {
			activePeersCounter++
		}
	}

	return activePeersCounter > 0
}

func (pd *PeerData) SendDataOnPeerId(data []byte, id peer.ID) bool {
	if err := pd.Peer.Connect(context.Background(), pd.Peer.Peerstore().PeerInfo(id)); err == nil {
		stream, err := pd.Peer.NewStream(context.Background(), id, PROTOCOL_ID)
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
