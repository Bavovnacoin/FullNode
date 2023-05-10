/*
	Tests how communication between two nodes performed by sending
	some data like blocks, txs, requesting time.
*/

package singleFunctionTesting

import (
	"bavovnacoin/ecdsa"
	"bavovnacoin/networking_p2p"
	"fmt"
	"math/big"
	"math/rand"
	"strings"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/multiformats/go-multiaddr"
)

type CommunicationTest struct {
	SingleFunctionTesting
	// TODO: value for ammount of peers
	peer1 networking_p2p.PeerData
	peer2 networking_p2p.PeerData
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

	return Peer, addr
}

func (ct *CommunicationTest) Launch() {
	var addr1 string
	var addr2 string

	ct.peer1.Peer, addr1 = ct.startPeer()
	ct.peer2.Peer, addr2 = ct.startPeer()

	ct.peer2.Peer, ct.peer2.OtherPeersIds, _ = ct.addOtherAddress(addr1, ct.peer2.OtherPeersIds, ct.peer2.Peer)
	ct.peer1.Peer, ct.peer1.OtherPeersIds, _ = ct.addOtherAddress(addr2, ct.peer1.OtherPeersIds, ct.peer1.Peer)

}
