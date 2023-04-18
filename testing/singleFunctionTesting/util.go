package singleFunctionTesting

import (
	"bavovnacoin/account"
	"fmt"
)

type SingleFunctionTesting struct{}

func (bvt *SingleFunctionTesting) genBlockTestAccounts(ammount int) {
	for i := 0; i < ammount; i++ {
		account.Wallet = append(account.Wallet, account.GenAccount(fmt.Sprint(len(account.Wallet))))
	}
}
