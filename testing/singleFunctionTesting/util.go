package singleFunctionTesting

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/byteArr"
	"bavovnacoin/hashing"
	"bavovnacoin/testing/account"
	"bavovnacoin/transaction"
	"fmt"
	"math/rand"
)

type SingleFunctionTesting struct{}

func (bvt *SingleFunctionTesting) genBlockTestAccounts(ammount int) {
	for i := 0; i < ammount; i++ {
		account.Wallet = append(account.Wallet, account.GenAccount(fmt.Sprint(len(account.Wallet))))
	}
}

func (bvt *SingleFunctionTesting) CreateTestTx(textForAddress string, random *rand.Rand) (transaction.Transaction, bool) {
	fee := random.Intn(5) + 1
	isGenLocktime := random.Intn(5)
	var locktime uint
	if isGenLocktime == 2 {
		locktime = uint(int(blockchain.BcLength+1) + random.Intn(2) + 1)
	}

	var outAddr []byteArr.ByteArr
	outAddr = append(outAddr, byteArr.ByteArr{})
	outAddr[0].SetFromHexString(hashing.SHA1(textForAddress), 20)

	var outValue []uint64
	outValue = append(outValue, 1000)

	tx, isValid := transaction.CreateTransaction(fmt.Sprint(0), outAddr, outValue, fee, locktime)

	return tx, isValid == ""
}
