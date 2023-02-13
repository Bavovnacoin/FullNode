package account

import (
	"bavovnacoin/byteArr"
	"bavovnacoin/cryption"
	"bavovnacoin/ecdsa"
	"bavovnacoin/hashing"
	"bavovnacoin/utxo"
	"fmt"
	"sort"
)

var CurrAccount Account

type Account struct {
	Id          string
	HashPass    string
	KeyPairList []ecdsa.KeyPair
	ArrId       int    `json:"-"`
	Balance     uint64 `json:"-"`
}

// Generates new account and set up a password to encode a private key
func GenAccount(password string) Account {
	ecdsa.InitValues()
	var newAcc Account

	newAcc.HashPass = hashing.SHA1(password)

	newKeyPair := ecdsa.GenKeyPair()
	newKeyPair.PrivKey = cryption.AES_encrypt(newKeyPair.PrivKey, password)

	newAcc.Id = fmt.Sprint(RightBoundAccNum + 1)
	RightBoundAccNum++
	newAcc.KeyPairList = append(newAcc.KeyPairList, newKeyPair)

	// Wallet = append(Wallet, newAcc)
	//WriteAccounts()
	return newAcc
}

func AddKeyPairToAccount(password string) string {
	if CurrAccount.HashPass == hashing.SHA1(password) {
		ecdsa.InitValues()
		newKeyPair := ecdsa.GenKeyPair()
		newKeyPair.PrivKey = cryption.AES_encrypt(newKeyPair.PrivKey, password)
		CurrAccount.KeyPairList = append(CurrAccount.KeyPairList, newKeyPair)
		Wallet[CurrAccount.ArrId] = CurrAccount
		//WriteAccounts()
	} else {
		return "Wrong password!"
	}
	return ""
}

func GetAccUtxo() []utxo.TXO {
	var accUtxo []utxo.TXO
	for i := 0; i < len(CurrAccount.KeyPairList); i++ {
		for j := 0; j < len(utxo.CoinDatabase); j++ {
			var currAccAddr byteArr.ByteArr
			currAccAddr.SetFromHexString(hashing.SHA1(CurrAccount.KeyPairList[i].PublKey), 20)
			if utxo.CoinDatabase[j].OutAddress.IsEqual(currAccAddr) {
				accUtxo = append(accUtxo, utxo.CoinDatabase[j])
			}
		}
	}
	sort.Slice(accUtxo, func(i, j int) bool {
		return accUtxo[i].Sum > accUtxo[j].Sum
	})
	return accUtxo
}

func GetBalHashOutInd(txHash byteArr.ByteArr, outInd int) uint64 {
	for j := 0; j < len(utxo.CoinDatabase); j++ {
		if txHash.IsEqual(utxo.CoinDatabase[j].OutTxHash) && utxo.CoinDatabase[j].TxOutInd == uint64(outInd) {
			return utxo.CoinDatabase[j].Sum
		}
	}
	return 0
}

func GetBalByAddress(address byteArr.ByteArr) uint64 {
	var sum uint64
	for i := 0; i < len(utxo.CoinDatabase); i++ {
		if address.IsEqual(utxo.CoinDatabase[i].OutAddress) {
			sum += utxo.CoinDatabase[i].Sum
		}
	}
	return sum
}

// A function counts all the UTXOs that is on specific public keys on user's account
func GetBalance() uint64 {
	CurrAccount.Balance = 0
	for i := 0; i < len(CurrAccount.KeyPairList); i++ {
		var address byteArr.ByteArr
		address.SetFromHexString(hashing.SHA1(CurrAccount.KeyPairList[i].PublKey), 20)
		CurrAccount.Balance += GetBalByAddress(address)
	}
	return CurrAccount.Balance
}

func PrintBalance() {
	GetBalance()
	var bal float64 = float64(CurrAccount.Balance) / 100000000.
	fmt.Printf("Balance: %.8f BVC\n", bal)
}

func getAccountInd(accountId string) int {
	for i := 0; i < len(Wallet); i++ {
		if Wallet[i].Id == accountId {
			Wallet[i].ArrId = i
			return i
		}
	}
	return -1
}

func InitAccount(accountId string) bool {
	ecdsa.InitValues()
	accInd := getAccountInd(accountId)
	if accInd != -1 {
		CurrAccount = Wallet[accInd]
		return true
	}
	return false
}

func SignData(hashMes string, kpInd int, pass string) (string, bool) {
	if CurrAccount.HashPass != hashing.SHA1(pass) {
		return "", true
	}
	kp := CurrAccount.KeyPairList[kpInd]
	kp.PrivKey = cryption.AES_decrypt(kp.PrivKey, pass)

	return ecdsa.Sign(hashMes, kp.PrivKey), false
}

func VerifData(hashMes string, kpInd int, signature string) bool {
	kp := CurrAccount.KeyPairList[kpInd]
	return ecdsa.Verify(kp.PublKey, signature, hashMes)
}
