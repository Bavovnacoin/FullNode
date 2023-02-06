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
	ScriptSig string
	OutInd    int
}

type Output struct {
	Address byteArr.ByteArr
	Sum     uint64
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
		message += fmt.Sprint(tx.Outputs[i].Sum)
	}
	return message
}

func genTxScriptSignatures(keyPair []ecdsa.KeyPair, passKey string, tx Transaction) Transaction {
	message := hashing.SHA1(GetCatTxFields(tx))
	// Signing message
	for i := 0; i < len(keyPair); i++ {
		tx.Inputs[i].ScriptSig = keyPair[i].PublKey + ecdsa.Sign(message, cryption.AES_decrypt(keyPair[i].PrivKey, passKey))
	}

	return tx
}

func ComputeTxSize(tx Transaction) int {
	size := 0
	size += 8 // 4 bytes for Version, 4 for locktime
	for i := 0; i < len(tx.Inputs); i++ {
		size += len(tx.Inputs[i].ScriptSig)
		size += 4 // Input out index size
		size += 20
		size += len(tx.Inputs[i].ScriptSig)
	}

	for i := 0; i < len(tx.Outputs); i++ {
		size += 8  // Output value
		size += 20 //len(tx.Outputs[i].Address)
	}
	return size
}

func getNextInpIndex(addressInp byteArr.ByteArr, utxoInputs []utxo.TXO, utxoInd int) int {
	ind := -1
	for i := 0; i <= utxoInd; i++ {
		if utxoInputs[i].OutAddress.IsEqual(addressInp) {
			ind++
		}
	}
	return ind
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

// TODO: change function according to new utxo structure
func GetTransInputs(sum uint64, accUtxo []utxo.TXO) ([]UtxoForInput, []utxo.TXO, uint64) {
	if accUtxo == nil {
		accUtxo = account.GetAccUtxo()
	}

	accUtxo = append(accUtxo, utxo.TXO{}) // Stub value for searching
	var utxoInput []UtxoForInput
	tempSum := uint64(0)

	//getNextInpIndex(accUtxo[0].OutAddress, accUtxo, 0)
	if len(accUtxo) == 1 && accUtxo[0].Sum >= sum {
		return append(utxoInput, UtxoForInput{accUtxo[0].OutTxHash, int(accUtxo[0].TxOutInd)}),
			accUtxo, accUtxo[0].Sum
	}

	for i := 1; i < len(accUtxo); i++ {
		if accUtxo[i-1].Sum >= sum-tempSum {
			if sum-tempSum > accUtxo[i].Sum {
				utxoInput = append(utxoInput, UtxoForInput{TxHash: accUtxo[i-1].OutTxHash, OutInd: int(accUtxo[i-1].TxOutInd)})
				return utxoInput, accUtxo, accUtxo[i-1].Sum + tempSum
			} else {
				continue
			}
		}
		utxoInput = append(utxoInput, UtxoForInput{accUtxo[i-1].OutTxHash,
			int(accUtxo[i-1].TxOutInd)})
		tempSum += accUtxo[i-1].Sum
	}
	return nil, accUtxo, tempSum
}

// Creates transaction
func CreateTransaction(passKey string, outAdr []byteArr.ByteArr, outSumVals []uint64,
	feePerByte int, locktime uint) (Transaction, string) {
	var tx Transaction
	tx.Locktime = locktime
	txSize := 0
	tx.Version = 0
	genSum := uint64(0)
	for i := 0; i < len(outSumVals); i++ {
		genSum += outSumVals[i]
	}

	// Genereting outputs
	var output []Output
	for i := 0; i < len(outAdr); i++ {
		var outVal Output
		outVal.Address = outAdr[i]
		outVal.Sum = uint64(outSumVals[i])
		output = append(output, outVal)
	}

	// Genereting inputs (including tx fee)
	var input []Input
	kpAcc := make([]ecdsa.KeyPair, len(account.CurrAccount.KeyPairList))
	copy(kpAcc, account.CurrAccount.KeyPairList)
	outTxSum := uint64(0)
	needSum := genSum + uint64(txSize)*uint64(feePerByte)

	var kpForSign []ecdsa.KeyPair
	for outTxSum < needSum { // Looking for tx fee
		kpForSign = []ecdsa.KeyPair{}
		inputs, _, outSum := GetTransInputs(needSum, nil)

		if needSum > outSum {
			return tx, "Not enough coins for payment. You need: " + fmt.Sprint(needSum) + ", you have: " + fmt.Sprint(account.GetBalance())
		}

		outTxSum = outSum
		for i := 0; i < len(inputs); i++ {
			var inpVal Input
			inpVal.TxHash = inputs[i].TxHash
			inpVal.OutInd = inputs[i].OutInd

			// Get private and public key for scriptSig generation
			for j := 0; j < len(kpAcc); j++ {
				var newAddr byteArr.ByteArr
				newAddr.SetFromHexString(hashing.SHA1(kpAcc[j].PublKey), 20)
				if newAddr.IsEqual(inputs[i].TxHash) {
					kpForSign = append(kpForSign, ecdsa.KeyPair{PrivKey: kpAcc[j].PrivKey, PublKey: kpAcc[j].PublKey})
				}
			}
			input = append(input, inpVal)
		}
		tx.Inputs = input
		tx.Outputs = output
		txSize = ComputeTxSize(tx)
		needSum = genSum + uint64(txSize)*uint64(feePerByte)
	}

	//Generating change output
	if outTxSum > genSum+uint64(txSize)*uint64(feePerByte) {
		account.AddKeyPairToAccount(passKey) // generate new keypair for the change
		kpLen := len(account.CurrAccount.KeyPairList)
		tx.Outputs = append(tx.Outputs, Output{Sum: uint64(outTxSum - (genSum + uint64(txSize)*uint64(feePerByte)))})
		tx.Outputs[len(tx.Outputs)-1].Address.SetFromHexString(hashing.SHA1(account.CurrAccount.KeyPairList[kpLen-1].PublKey), 20)
	}
	tx = genTxScriptSignatures(kpForSign, passKey, tx)
	return tx, ""
}

func GetInputSum(inp []Input) uint64 {
	var sum uint64 = 0
	for i := 0; i < len(inp); i++ {
		sum += account.GetBalHashOutInd(inp[i].TxHash, inp[i].OutInd)
	}
	return sum
}

func GetOutputSum(out []Output) uint64 {
	var sum uint64 = 0
	for i := 0; i < len(out); i++ {
		sum += out[i].Sum
	}
	return sum
}

func GetTxFee(tx Transaction) uint64 {
	return GetInputSum(tx.Inputs) - GetOutputSum(tx.Outputs)
}

/*
Just to show that everything works fine.

Some information is not stored in the transaction structure,
but received in this function.
*/
func PrintTransaction(tx Transaction) {
	fmt.Printf("Version: %d\nLocktime %d\n", tx.Version, tx.Locktime)
	println("Inputs:")
	var inpSum uint64
	for i := 0; i < len(tx.Inputs); i++ {
		curVal := account.GetBalHashOutInd(tx.Inputs[i].TxHash, tx.Inputs[i].OutInd)
		inpSum += curVal
		fmt.Printf("%d. HashAddress: %s (Bal: %d)\nOut index: %d\nScriptSig: %s\n", i, tx.Inputs[i].TxHash.ToHexString(), curVal,
			tx.Inputs[i].OutInd, tx.Inputs[i].ScriptSig)
	}
	println("\nOutputs:")
	var outSum uint64
	for i := 0; i < len(tx.Outputs); i++ {
		outSum += tx.Outputs[i].Sum
		fmt.Printf("%d. HashAddress: %s. Sum: %d\n", i, tx.Outputs[i].Address.ToHexString(), tx.Outputs[i].Sum)
	}
	print("(Fee: ")
	println(inpSum-outSum, ")")
}

// Verifies transaction
func VerifyTransaction(tx Transaction) bool {
	if tx.Version == 0 {
		var inpSum uint64
		var outSum uint64
		hashMesOfTx := hashing.SHA1(GetCatTxFields(tx))
		var txoArr []utxo.TXO

		// Checking signatures and unique inputs
		for i := 0; i < len(tx.Inputs); i++ {
			if len(tx.Inputs[i].ScriptSig) < 66 {
				return false
			}

			pubKey := tx.Inputs[i].ScriptSig[:66]
			sign := tx.Inputs[i].ScriptSig[66:]
			if !ecdsa.Verify(pubKey, sign, hashMesOfTx) {
				return false
			}
			curVal := account.GetBalHashOutInd(tx.Inputs[i].TxHash, tx.Inputs[i].OutInd)
			inpSum += curVal

			isExists := utxo.IsUtxoExists(tx.Inputs[i].TxHash, uint64(tx.Inputs[i].OutInd))
			if isExists {
				outSum += txoArr[i].Sum
			} else {
				return false
			}
		}

		// Checking presence of coins to be spent
		if inpSum < outSum {
			return false
		}
	}
	return true
}
