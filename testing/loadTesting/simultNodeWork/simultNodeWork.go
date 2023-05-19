package simultNodeWork

import (
	"bavovnacoin/hashing"
	"bavovnacoin/node/node_settings"
	"bavovnacoin/testing"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

type SimultaneousNodeWork struct {
	NodesAmmount        int
	consoles            []*exec.Cmd
	blocksToMineAmmount int

	addresses [][]string

	pathToNodes string
	rootPath    string

	peerConnAmmount int

	source rand.Source
	random *rand.Rand
}

func (snw *SimultaneousNodeWork) ExecCommand(name string, params ...string) (*exec.Cmd, string) {
	cmd := exec.Command(name, strings.Join(params, " "))
	output, err := cmd.Output()
	if err != nil {
		return cmd, err.Error()
	}
	return cmd, string(output)
}

func (snw *SimultaneousNodeWork) CreateSettings(nodeId int) {
	var settings node_settings.NodeSettings
	settings.TxMinFee = 0
	settings.MiningThreads = 0
	settings.NodeType = 0
	settings.MyAddress = testing.GetRandAddress("/ip4/127.0.0.1/tcp/")

	nodeIdStr := fmt.Sprint(nodeId + 1)
	settings.SetPrivKey("password")

	settings.HashPass.SetFromHexString(hashing.SHA1("password"), 20)

	settings.RewardAddress = strings.Repeat(nodeIdStr, 40/len(nodeIdStr))
	settings.RewardAddress += nodeIdStr[:40-len(settings.RewardAddress)]

	settings.RPCip = testing.GetRandAddress("localhost:")
	settings.WriteSettings()

	settings.DecryptPrivKey("password")
	node_settings.Settings.PrivKeyDecrypted = settings.PrivKeyDecrypted
	node_settings.Settings.MyAddress = settings.MyAddress

	var peer_rpc_addresses []string
	peer_rpc_addresses = append(peer_rpc_addresses, snw.GetPeerAddress(settings.GetPrivKey()))
	peer_rpc_addresses = append(peer_rpc_addresses, settings.RPCip)
	snw.addresses = append(snw.addresses, peer_rpc_addresses)
}

func (snw *SimultaneousNodeWork) CreateTestNodes() {
	println("Creating test nodes...")
	snw.ExecCommand("cmd", "/C", "mkdir", snw.pathToNodes+"node"+fmt.Sprint(0))
	snw.ExecCommand("cmd", "/C", "go", "build", "-o", snw.pathToNodes+"node"+fmt.Sprint(0), "main.go")

	for i := 0; i < snw.NodesAmmount; i++ {
		snw.ExecCommand("cmd", "/C", "mkdir", snw.pathToNodes+"node"+fmt.Sprint(i))
		cmd, _ := snw.ExecCommand("cmd", "/C", "copy", snw.pathToNodes+"node0\\main.exe", snw.pathToNodes+"node"+fmt.Sprint(i))
		snw.consoles = append(snw.consoles, cmd)
	}

	println("Nodes are created")
}

func (snw *SimultaneousNodeWork) InitSettingsFiles() {
	for i := 0; i < len(snw.consoles); i++ {
		os.Chdir(snw.rootPath + "\\" + snw.pathToNodes + "node" + fmt.Sprint(i))
		snw.CreateSettings(i)
	}
}

func (snw *SimultaneousNodeWork) GetPeerAddress(privKey []byte) string {
	privateKey, err := crypto.UnmarshalSecp256k1PrivateKey(privKey)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	publicKey := privateKey.GetPublic()
	peerID, _ := peer.IDFromPublicKey(publicKey)

	return node_settings.Settings.MyAddress + "/" + peerID.Pretty()
}

func (snw *SimultaneousNodeWork) ManageConnections() {
	currAddr := snw.GetPeerAddress(node_settings.Settings.GetPrivKey())
	addressesLocal := [][]string{}

	for i := 0; i < len(snw.addresses); i++ {
		if snw.addresses[i][0] != currAddr {
			addressesLocal = append(addressesLocal, snw.addresses[i])
		}
	}

	for i := 0; i < snw.peerConnAmmount && len(addressesLocal) != 0; i++ {
		ind := snw.random.Int() % len(addressesLocal)
		node_settings.Settings.OtherNodesAddresses = append(node_settings.Settings.OtherNodesAddresses, addressesLocal[ind])
		addressesLocal = append(addressesLocal[:ind], addressesLocal[ind+1:]...)
	}

	node_settings.Settings.WriteSettings()
}

func (snw *SimultaneousNodeWork) LaunchTestNodes() {
	for i := 0; i < len(snw.consoles); i++ {
		os.Chdir(snw.rootPath + "\\" + snw.pathToNodes + "node" + fmt.Sprint(i))
		cwd, _ := os.Getwd()

		node_settings.Settings.GetSettings()
		node_settings.Settings.DecryptPrivKey("password")

		snw.ManageConnections()

		go snw.ExecCommand("cmd", "/c", "start", cwd+"\\main.exe -snw -l")
		//time.Sleep(100 * time.Millisecond)
		time.Sleep(4 * time.Second)
	}
	time.Sleep(1 * time.Second)
}

func (snw *SimultaneousNodeWork) Launch(nodesAmmount int) {
	snw.NodesAmmount = nodesAmmount
	snw.peerConnAmmount = 2
	snw.source = rand.NewSource(time.Now().Unix())
	snw.random = rand.New(snw.source)
	snw.blocksToMineAmmount = 10

	snw.rootPath, _ = os.Getwd()
	snw.pathToNodes = "data\\SimNodeWorkTest\\"

	snw.CreateTestNodes()
	snw.InitSettingsFiles()
	//defer os.RemoveAll(snw.pathToNodes)

	//blockchain.STARTBITS
	//snw.LaunchTestNodes()
}
