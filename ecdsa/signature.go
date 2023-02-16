package ecdsa

import (
	"math/big"
	"math/rand"
	"time"
)

func Sign(hashMes string, privKey string) string { // invalid pr key starts with 0 - an empty byte!b
	bigHash, isHashValid := hexToBig(hashMes)
	prKNum, isKeyValid := hexToBig(privKey)

	if !isHashValid || !isKeyValid {
		return ""
	}

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	randK := big.NewInt(0).Rand(r, n)
	x := big.NewInt(0).Mod(multiply(randK, g).x, n)

	hashKeyX := big.NewInt(0).Add(big.NewInt(1).Mul(prKNum, x), bigHash)
	y := big.NewInt(0).Mod(big.NewInt(1).Mul(inverse(randK, n), hashKeyX), n)
	return fillZero(x.String(), 78) + fillZero(y.String(), 78)
}

func Verify(pubKey, sign, hashMes string) bool {
	pKey := decompressPubKey(pubKey)

	if len(sign) < 78 {
		return false
	}

	var signPoint Point
	errX := true
	errY := true
	signPoint.x, errX = big.NewInt(0).SetString(sign[:78], 10)
	signPoint.y, errY = big.NewInt(0).SetString(sign[78:], 10)
	if !errX || !errY {
		return false
	}

	bigHash, isHashValid := hexToBig(hashMes)
	if !isHashValid {
		return false
	}

	inv := inverse(signPoint.y, n)
	u1 := big.NewInt(1).Mod(big.NewInt(1).Mul(bigHash, inv), n)
	u2 := big.NewInt(1).Mod(big.NewInt(1).Mul(signPoint.x, inv), n)
	p3 := add(multiply(u1, g), multiply(u2, pKey))

	if p3.x.String() == signPoint.x.String() {
		return true
	}

	return false
}
