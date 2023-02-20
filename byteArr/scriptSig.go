package byteArr

type ScriptSig struct {
	ByteArr
}

func (scriptSig ScriptSig) GetPubKey() ByteArr {
	var pubKey ByteArr
	pubKey.ByteArr = scriptSig.ByteArr.ByteArr[:33]
	return pubKey
}

func (scriptSig ScriptSig) GetSignature() ByteArr {
	var sign ByteArr
	sign.ByteArr = scriptSig.ByteArr.ByteArr[33:]
	return sign
}
