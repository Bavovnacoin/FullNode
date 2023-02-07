package byteArr

type ScriptSig struct {
	ByteArr
}

func (scriptSig ScriptSig) GetPubKey() ByteArr {
	var pubKey ByteArr
	pubKey.byteArr = scriptSig.ByteArr.byteArr[:33]
	return pubKey
}

func (scriptSig ScriptSig) GetSignature() ByteArr {
	var sign ByteArr
	sign.byteArr = scriptSig.ByteArr.byteArr[33:]
	return sign
}
