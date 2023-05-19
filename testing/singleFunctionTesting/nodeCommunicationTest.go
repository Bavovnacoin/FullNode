/*
	Tests if a communication between two peers works fine
	by requesting local node time
*/

package singleFunctionTesting

import (
	"bavovnacoin/ecdsa"
	"bavovnacoin/networking_p2p"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/multiformats/go-multiaddr"
)

type CommunicationTest struct {
	SingleFunctionTesting

	peer1 networking_p2p.PeerData
	peer2 networking_p2p.PeerData

	mainPeer         networking_p2p.PeerData
	otherPeers       []networking_p2p.PeerData
	otherPeerAmmount int

	source rand.Source
	random *rand.Rand
}

func (ct *CommunicationTest) getRandAddress() string {
	return "/ip4/127.0.0.1/tcp/" + fmt.Sprint(rand.Intn(9000)+1000)
}

func (ct *CommunicationTest) addOtherAddress(address string, OtherPeerIds []peer.ID, Peer host.Host) (host.Host, []peer.ID, bool) {
	arr := strings.Split(address, "/")

	maddr, err := multiaddr.NewMultiaddr(strings.Join(arr[:len(arr)-1], "/"))
	if err != nil {
		return Peer, OtherPeerIds, false
	}

	id, res := peer.Decode(arr[len(arr)-1])
	if res != nil {
		return Peer, OtherPeerIds, false
	}

	OtherPeerIds = append(OtherPeerIds, id)
	Peer.Peerstore().AddAddrs(id, []multiaddr.Multiaddr{maddr}, peerstore.PermanentAddrTTL)
	return Peer, OtherPeerIds, true
}

func (ct *CommunicationTest) startPeer() (host.Host, string) {
	ecdsa.InitValues()
	pk, _ := new(big.Int).SetString(ecdsa.GenPrivKey(), 16)
	privKey, _ := crypto.UnmarshalSecp256k1PrivateKey(pk.Bytes())

	addr := ct.getRandAddress()

	var err error
	Peer, err := libp2p.New(
		libp2p.Identity(privKey),
		libp2p.ListenAddrStrings(addr),
	)

	if err == nil {
		Peer.SetStreamHandler(networking_p2p.PROTOCOL_ID, networking_p2p.Peer.StreamHandler)
		fmt.Printf("Started peer on address %s/%s\n", Peer.Addrs()[0], Peer.ID().Pretty())
	} else {
		fmt.Println("Unable to start a peer.", err)
	}

	return Peer, fmt.Sprintf("%s/%s", Peer.Addrs()[0], Peer.ID().Pretty())
}

func (ct *CommunicationTest) StartOtherPeers() {
	for i := 0; i < ct.otherPeerAmmount; i++ {
		peer, addr := ct.startPeer()
		var peerData networking_p2p.PeerData
		peerData.Peer = peer
		ct.otherPeers = append(ct.otherPeers, peerData)

		ct.mainPeer.Peer, ct.mainPeer.OtherPeersIds, _ = ct.addOtherAddress(addr, ct.mainPeer.OtherPeersIds, ct.mainPeer.Peer)
	}

	networking_p2p.Peer.Peer = ct.otherPeers[0].Peer
}

func (ct *CommunicationTest) TestCommWithTimeRequest() {
	ct.mainPeer.RequestNodesTime("")
}

func (ct *CommunicationTest) PrintResult() {
	println("Results:")
	fmt.Printf("Got %d/%d responces\n", len(networking_p2p.NodesTime), ct.otherPeerAmmount)

	if len(networking_p2p.NodesTime) == ct.otherPeerAmmount {
		println("Test is passed")
	} else {
		println("Test is not passed")
	}
}

func (ct *CommunicationTest) Launch(otherPeerAmmount int) {
	ct.source = rand.NewSource(time.Now().Unix())
	ct.random = rand.New(ct.source)
	ct.otherPeerAmmount = otherPeerAmmount

	println("Starting main peer")
	ct.mainPeer.Peer, _ = ct.startPeer()
	networking_p2p.Peer.Peer = ct.mainPeer.Peer

	println("Starting other peers")
	ct.StartOtherPeers()

	ct.TestCommWithTimeRequest()
	time.Sleep(1 * time.Second)

	ct.PrintResult()
}
