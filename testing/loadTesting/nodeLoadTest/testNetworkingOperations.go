package nodeLoadTest

import (
	"bavovnacoin/byteArr"
	"bavovnacoin/ecdsa"
	"bavovnacoin/networking"
	"bavovnacoin/networking_p2p"
	"bavovnacoin/node/node_settings"
	"bavovnacoin/transaction"
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"net/rpc"
	"strings"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/multiformats/go-multiaddr"
)

type Connection struct {
	client *rpc.Client
}

type Reply struct {
	Data []byte
}

func getRandAddress() string {
	return "/ip4/127.0.0.1/tcp/" + fmt.Sprint(rand.Intn(9000)+1000)
}

func InitTestSettings() {
	node_settings.Settings.MyAddress = getRandAddress()
	ecdsa.InitValues()
	node_settings.Settings.PrivKeyDecrypted, _ = hex.DecodeString(ecdsa.GenPrivKey())
}

func startTestPeer() (host.Host, string) {
	ecdsa.InitValues()
	pk, _ := new(big.Int).SetString(ecdsa.GenPrivKey(), 16)
	privKey, _ := crypto.UnmarshalSecp256k1PrivateKey(pk.Bytes())

	addr := getRandAddress()

	var err error
	Peer, err := libp2p.New(
		libp2p.Identity(privKey),
		libp2p.ListenAddrStrings(addr),
	)

	if err == nil {
		Peer.SetStreamHandler(networking_p2p.PROTOCOL_ID, networking_p2p.Peer.StreamHandler)
	}

	return Peer, fmt.Sprintf("%s/%s", Peer.Addrs()[0], Peer.ID().Pretty())
}

func addOtherAddress(address string, pr networking_p2p.PeerData) (networking_p2p.PeerData, bool) {
	arr := strings.Split(address, "/")

	maddr, err := multiaddr.NewMultiaddr(strings.Join(arr[:len(arr)-1], "/"))
	if err != nil {
		return pr, false
	}

	id, res := peer.Decode(arr[len(arr)-1])
	if res != nil {
		return pr, false
	}

	pr.Peer.Peerstore().AddAddrs(id, []multiaddr.Multiaddr{maddr}, peerstore.PermanentAddrTTL)
	return pr, true
}

func getIpFromNodeAddr(nodeAddr string) string {
	ipArr := strings.Split(nodeAddr, "/")
	return ipArr[2] + ":" + ipArr[3]
}

func (c *Connection) Establish() bool {
	var err error
	node_settings.Settings.GetSettings()
	c.client, err = rpc.Dial("tcp", "localhost:8080") //node_settings.Settings.MyAddress
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (c *Connection) Close() {
	c.client.Close()
}

func (c *Connection) ToByteArr(data any) ([]byte, bool) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(data)
	if err != nil {
		return nil, false
	}

	return buffer.Bytes(), true
}

func (c *Connection) FromByteArr(dataByte []byte, data interface{}) bool {
	buf := bytes.NewBuffer(dataByte)
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(data)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func (c *Connection) SendTransaction(tx transaction.Transaction, isAccepted *bool) bool {
	byteArr, isConv := c.ToByteArr(tx)
	if !isConv {
		return false
	}

	var l networking.Listener
	var repl networking.Reply
	err := l.AddNewTxToMemp(byteArr, &repl)
	if err != nil {
		log.Println(err)
		return false
	}

	if repl.Data[0] == 1 {
		*isAccepted = true
	} else {
		*isAccepted = false
	}
	return true
}

func (c *Connection) IsAddrExist(addr byteArr.ByteArr) bool {
	var repl Reply
	err := c.client.Call("Listener.IsAddrExist", addr.ByteArr, &repl)
	if err != nil {
		log.Println(err)
		return false
	}

	if repl.Data[0] == 0 {
		return false
	}

	return true
}

// GetMyUtxo - in wallet
func (c *Connection) GetUtxoByAddress(addresses []byteArr.ByteArr) bool {
	byteArr, isConv := c.ToByteArr(addresses)
	if !isConv {
		return false
	}

	var repl networking.Reply
	var l networking.Listener
	err := l.GetUtxoByAddr(byteArr, &repl)
	if err != nil {
		log.Println(err)
		return false
	}

	return true

}
