package transaction

import (
	"bavovnacoin/account"
	"bavovnacoin/byteArr"
	"bavovnacoin/cryption"
	"bavovnacoin/ecdsa"
	"bavovnacoin/hashing"
	"bavovnacoin/utxo"
	"fmt"
)

type Input struct {
	TxHash    byteArr.ByteArr
	ScriptSig byteArr.ScriptSig
	OutInd    int
}

type Output struct {
	Address byteArr.ByteArr
	Value   uint64
}

type Transaction struct {
	Version  uint
	Locktime uint
	Inputs   []Input
	Outputs  []Output
}

type UtxoForInput struct {
	TxHash byteArr.ByteArr
	OutInd int
}

// Generating message for signing (SCRIPTHASH_ALL)
func GetCatTxFields(tx Transaction) string {
	message := ""
	message += fmt.Sprint(tx.Version)
	message += fmt.Sprint(tx.Locktime)
	for i := 0; i < len(tx.Inputs); i++ {
		message += tx.Inputs[i].TxHash.ToHexString()
		message += fmt.Sprint(tx.Inputs[i].OutInd)
	}
	for i := 0; i < len(tx.Outputs); i++ {
		message += tx.Outputs[i].Address.ToHexString()
		message += fmt.Sprint(tx.Outputs[i].Value)
	}
	return message
}

func genTxScriptSignatures(keyPair []ecdsa.KeyPair, passKey string, tx Transaction) Transaction {
	message := hashing.SHA1(GetCatTxFields(tx))
	// Signing message
	for i := 0; i < len(keyPair); i++ {
		//tx.Inputs[i].ScriptSig = keyPair[i].PublKey + ecdsa.Sign(message, cryption.AES_decrypt(keyPair[i].PrivKey, passKey))
		tx.Inputs[i].ScriptSig.SetFromHexString(keyPair[i].PublKey+ecdsa.Sign(message, cryption.AES_decrypt(keyPair[i].PrivKey, passKey)), 111)
	}

	return tx
}

func ComputeTxSize(tx Transaction) int {
	size := 0
	size += 8 // 4 bytes for Version, 4 for locktime
	for i := 0; i < len(tx.Inputs); i++ {
		size += 111 //len(tx.Inputs[i].ScriptSig)
		size += 4   // Input out index size
		size += 20
		//size += len(tx.Inputs[i].ScriptSig)
	}

	for i := 0; i < len(tx.Outputs); i++ {
		size += 8  // Output value
		size += 20 //len(tx.Outputs[i].Address)
	}
	return size
}

/*
Algorithm of effective transaction inputs search:
iterate utxo of a specific account and check two neighboring values.
At the beginning we add a stub UTXO (it will not be added to the database).
We are looking for a minimum value (checking left neighbor)
that is higher or equal to the required sum (minus sum that we have already found).
If a right neighbor is less than needed sum, we keep iterating, because there is a chance
of finding better variant.
*/
func GetTransInputs(value uint64, accUtxo []utxo.TXO) ([]UtxoForInput, []utxo.TXO, uint64) {
	if accUtxo == nil {
		accUtxo = account.GetAccUtxo()
	}

	accUtxo = append(accUtxo, utxo.TXO{}) // Stub value for searching
	var utxoInput []UtxoForInput
	tempValue := uint64(0)

	if len(accUtxo) == 1 && accUtxo[0].Value >= value {
		return append(utxoInput, UtxoForInput{accUtxo[0].OutTxHash, int(accUtxo[0].TxOutInd)}),
			accUtxo, accUtxo[0].Value
	}

	for i := 1; i < len(accUtxo); i++ {
		if accUtxo[i-1].Value >= value-tempValue && accUtxo[i-1].Value != 0 {
			if value-tempValue > accUtxo[i].Value {
				utxoInput = append(utxoInput, UtxoForInput{TxHash: accUtxo[i-1].OutTxHash, OutInd: int(accUtxo[i-1].TxOutInd)})
				return utxoInput, accUtxo, accUtxo[i-1].Value + tempValue
			} else {
				continue
			}
		}
		utxoInput = append(utxoInput, UtxoForInput{accUtxo[i-1].OutTxHash,
			int(accUtxo[i-1].TxOutInd)})
		tempValue += accUtxo[i-1].Value
	}
	return nil, accUtxo, tempValue
}

// Creates transaction
func CreateTransaction(passKey string, outAdr []byteArr.ByteArr, outVals []uint64,
	feePerByte int, locktime uint) (Transaction, string) {
	var tx Transaction
	tx.Locktime = locktime
	txSize := 0
	tx.Version = 0
	genValue := uint64(0)
	for i := 0; i < len(outVals); i++ {
		genValue += outVals[i]
	}

	// Genereting outputs
	var output []Output
	for i := 0; i < len(outAdr); i++ {
		var outVal Output
		outVal.Address = outAdr[i]
		outVal.Value = uint64(outVals[i])
		output = append(output, outVal)
	}

	// Genereting inputs (including tx fee)
	var input []Input
	kpAcc := make([]ecdsa.KeyPair, len(account.CurrAccount.KeyPairList))
	copy(kpAcc, account.CurrAccount.KeyPairList)
	outTxValue := uint64(0)
	needValue := genValue + uint64(txSize)*uint64(feePerByte)

	var kpForSign []ecdsa.KeyPair
	for outTxValue < needValue { // Looking for tx fee
		kpForSign = []ecdsa.KeyPair{}
		inputs, txo, outInpValue := GetTransInputs(needValue, nil)
		if needValue > outInpValue {
			return tx, "Not enough coins for payment. You need: " + fmt.Sprint(needValue) + ", you have: " + fmt.Sprint(account.GetBalance())
		}

		outTxValue = outInpValue
		for i := 0; i < len(inputs); i++ {
			var inpVal Input
			inpVal.TxHash = inputs[i].TxHash
			inpVal.OutInd = inputs[i].OutInd

			// Get private and public key for scriptSig generation
			isFound := false
			for j := 0; j < len(kpAcc); j++ {
				var newAddr byteArr.ByteArr
				newAddr.SetFromHexString(hashing.SHA1(kpAcc[j].PublKey), 20)
				for k := 0; k < len(txo); k++ {
					if inputs[i].OutInd == int(txo[k].TxOutInd) && newAddr.IsEqual(txo[k].OutAddress) {
						kpForSign = append(kpForSign, ecdsa.KeyPair{PrivKey: kpAcc[j].PrivKey, PublKey: kpAcc[j].PublKey})
						isFound = true
						break
					}
				}
				if isFound {
					break
				}
			}
			input = append(input, inpVal)
		}
		tx.Inputs = input
		tx.Outputs = output
		txSize = ComputeTxSize(tx)
		needValue = genValue + uint64(txSize)*uint64(feePerByte)
	}

	//Generating change output
	if outTxValue > genValue+uint64(txSize)*uint64(feePerByte) {
		account.AddKeyPairToAccount(passKey) // generate new keypair for the change
		kpLen := len(account.CurrAccount.KeyPairList)
		tx.Outputs = append(tx.Outputs, Output{Value: uint64(outTxValue - (genValue + uint64(txSize)*uint64(feePerByte)))})
		tx.Outputs[len(tx.Outputs)-1].Address.SetFromHexString(hashing.SHA1(account.CurrAccount.KeyPairList[kpLen-1].PublKey), 20)
	}

	tx = genTxScriptSignatures(kpForSign, passKey, tx)
	return tx, ""
}

func GetInputValue(inp []Input) uint64 {
	var value uint64 = 0
	for i := 0; i < len(inp); i++ {
		value += account.GetBalHashOutInd(inp[i].TxHash, inp[i].OutInd)
	}
	return value
}

func GetOutputValue(out []Output) uint64 {
	var value uint64 = 0
	for i := 0; i < len(out); i++ {
		value += out[i].Value
	}
	return value
}

func GetTxFee(tx Transaction) uint64 {
	return GetInputValue(tx.Inputs) - GetOutputValue(tx.Outputs)
}

/*
Just to show that everything works fine.

Some information is not stored in the transaction structure,
but received in this function.
*/
func PrintTransaction(tx Transaction) {
	fmt.Printf("Version: %d\nLocktime %d\n", tx.Version, tx.Locktime)
	println("Inputs:")
	var inpValue uint64
	for i := 0; i < len(tx.Inputs); i++ {
		curVal := account.GetBalHashOutInd(tx.Inputs[i].TxHash, tx.Inputs[i].OutInd)
		inpValue += curVal
		fmt.Printf("%d. HashAddress: %s (Bal: %d)\nOut index: %d\nScriptSig: %s\n", i, tx.Inputs[i].TxHash.ToHexString(), curVal,
			tx.Inputs[i].OutInd, tx.Inputs[i].ScriptSig.ToHexString())
	}
	println("\nOutputs:")
	var outValue uint64
	for i := 0; i < len(tx.Outputs); i++ {
		outValue += tx.Outputs[i].Value
		fmt.Printf("%d. HashAddress: %s. Value: %d\n", i, tx.Outputs[i].Address.ToHexString(), tx.Outputs[i].Value)
	}
	print("(Fee: ")
	println(inpValue-outValue, ")")
}

// Verifies transaction
func VerifyTransaction(tx Transaction) bool {
	if tx.Version == 0 {
		var inpValue uint64
		var outValue uint64
		hashMesOfTx := hashing.SHA1(GetCatTxFields(tx))

		// Checking signatures and unique inputs
		for i := 0; i < len(tx.Inputs); i++ {
			if len(tx.Inputs[i].ScriptSig.ToHexString()) < 66 {
				return false
			}

			pubKey := tx.Inputs[i].ScriptSig.GetPubKey().ToHexString()
			sign := tx.Inputs[i].ScriptSig.GetSignature().ToHexString()
			if !ecdsa.Verify(pubKey, sign, hashMesOfTx) {
				return false
			}

			curVal := account.GetBalHashOutInd(tx.Inputs[i].TxHash, tx.Inputs[i].OutInd)
			inpValue += curVal
		}

		for i := 0; i < len(tx.Outputs); i++ {
			outValue += tx.Outputs[i].Value
		}

		// Checking presence of coins to be spent
		if inpValue < outValue {
			return false
		}
	}
	return true
}
