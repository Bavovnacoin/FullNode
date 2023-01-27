package test

import (
	"bavovnacoin/account"
	"bavovnacoin/blockchain"
	"bavovnacoin/ecdsa"
	"bufio"
	"fmt"
	"os"
)

var node_working bool = true
var network_accounts []account.Account

func commandHandler() {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	if text == "stop\r\n" {
		node_working = false
	}

}

func process() {
	blockchain.InitBlockchain()
	var genesisAccKeyPair []ecdsa.KeyPair
	genesisAccKeyPair = append(genesisAccKeyPair, ecdsa.KeyPair{PrivKey: "d966fded26f23d50bb1223cdc6efe4cfebc9f2d6967cb570122c040baf5d42091953a2ba6466963351a4c6bc616e1858de87de02724cc89d9306a62b6d29fab6",
		PublKey: "033361587c679cf9476949cb7cdd15c07d6f2f9674886333f667bfedb87635a4b4"})
	network_accounts = append(network_accounts, account.Account{Id: "0",
		HashPass: "b6589fc6ab0dc82cf12099d1c2d40ab994e8410c", KeyPairList: genesisAccKeyPair})

	for node_working {
		createAccoundRandom()
		tx, _ := createRandomTransaction()
		if tx.Inputs != nil {
			if blockchain.AddTxToMempool(tx) {
				println("Tx added to mempool")
				println(fmt.Sprint(len(blockchain.Mempool)) + " - mempool len")
			}
		}
		addBlock()
	}
}

func Test1() {
	go process()
	commandHandler()
}
