package testing

import (
	"bavovnacoin/account"
	"bavovnacoin/byteArr"
	"bavovnacoin/cryption"
	"bavovnacoin/ecdsa"
	"bavovnacoin/hashing"
	"bavovnacoin/transaction"
	"bavovnacoin/txo"
	"fmt"
	"math/rand"
	"time"
)

func GenTestAccounts(ammount int) {
	for i := 0; i < ammount; i++ {
		account.Wallet = append(account.Wallet, account.GenAccount(fmt.Sprint(i)))
	}
}

func GenUtxo(i int, random *rand.Rand) {
	var outAddr byteArr.ByteArr
	outAddr.SetFromHexString(hashing.SHA1(account.Wallet[i].KeyPairList[0].PublKey), 20)

	var outTxHash byteArr.ByteArr
	outTxHash.SetFromHexString(hashing.SHA1(fmt.Sprint(i)+fmt.Sprint(time.Now().Unix())), 20)

	newTxo := txo.TXO{OutTxHash: outTxHash, Value: (random.Uint64()%700000 + 300000), OutAddress: outAddr}
	txo.AddUtxo(newTxo.OutTxHash, newTxo.TxOutInd, newTxo.Value, newTxo.OutAddress, newTxo.BlockHeight)
}

func GenTestUtxo(ammount int, random *rand.Rand) {
	for i := 0; i < ammount; i++ {
		GenUtxo(i, random)
	}
}

func MakeTxIncorrect(tx transaction.Transaction, incTxCounter int, random *rand.Rand) (transaction.Transaction, string) {
	if incTxCounter%3 == 0 { // Wrong tx hash in input
		var incTxHash byteArr.ByteArr
		inpId := 0
		if len(tx.Inputs) != 0 {
			inpId = random.Intn(len(tx.Inputs))
		}
		tx.Inputs[inpId].TxHash = incTxHash

		return tx, "Wrong tx hash in input num. " + fmt.Sprint(inpId)
	} else if incTxCounter%3 == 1 { // Wrong tx signature in input
		var incTxSig byteArr.ScriptSig
		fakeSigAccId := 0
		inpId := 0
		if len(tx.Inputs) != 0 {
			inpId = random.Intn(len(tx.Inputs))
			fakeSigAccId = random.Intn(len(tx.Inputs))
		}

		fakeSigPrivKey := cryption.AES_decrypt(account.Wallet[fakeSigAccId].KeyPairList[0].PrivKey, fmt.Sprint(fakeSigAccId))
		fakeSig := ecdsa.Sign(hashing.SHA1("Glory to Ukraine"), fakeSigPrivKey)

		incTxSig.SetFromHexString(tx.Inputs[inpId].ScriptSig.GetPubKey().ToHexString()+fakeSig, 111)
		tx.Inputs[inpId].ScriptSig = incTxSig
		return tx, "Wrong signature in input num. " + fmt.Sprint(inpId)
	} else { // Wrong tx output value
		outId := 0
		if len(tx.Inputs) != 0 {
			outId = random.Intn(len(tx.Outputs))
		}

		tx.Outputs[outId].Value = ^uint64(0)
		return tx, "Too high output value in output num. " + fmt.Sprint(outId)
	}
}

func GenValidTx(currAccId int, randLocktime int, random *rand.Rand) transaction.Transaction {
	account.CurrAccount = account.Wallet[currAccId]
	account.GetBalance()
	var outAddr byteArr.ByteArr
	outAddr.SetFromHexString(hashing.SHA1(account.Wallet[currAccId].KeyPairList[0].PublKey), 20)

	var outAddrTx []byteArr.ByteArr
	outAddrTx = append(outAddrTx, outAddr)

	var outValTx []uint64
	outValTx = append(outValTx, uint64(txo.CoinDatabase[currAccId].Value/((random.Uint64()%10)+3)))

	newTx, _ := transaction.CreateTransaction(fmt.Sprint(currAccId), outAddrTx, outValTx, random.Intn(10), uint(random.Intn(randLocktime)))
	return newTx
}

func GenRandTxs(txAmmount, incorrectTxAmmount int, random *rand.Rand) ([]transaction.Transaction, []string) {
	var randTxs []transaction.Transaction
	var txIncorrMessages []string

	var step int = int(txAmmount / incorrectTxAmmount)
	var incTxInd int = -1
	var incTxCounter int

	if incorrectTxAmmount != 0 {
		stStep := step * incTxCounter
		incTxInd = random.Intn(step) + stStep

		incTxCounter++
	}

	for i := 0; i < txAmmount; i++ {
		newTx := GenValidTx(i, 3, random)

		if i == incTxInd && incTxCounter <= incorrectTxAmmount {
			stStep := step * incTxCounter
			if incorrectTxAmmount-1 == incTxCounter {
				step = txAmmount - stStep
			}
			incTxInd = random.Intn(step) + stStep

			var message string
			newTx, message = MakeTxIncorrect(newTx, incTxCounter, random)
			txIncorrMessages = append(txIncorrMessages, message)
			incTxCounter++
		} else {
			txIncorrMessages = append(txIncorrMessages, "")
		}

		randTxs = append(randTxs, newTx)
	}

	return randTxs, txIncorrMessages
}
