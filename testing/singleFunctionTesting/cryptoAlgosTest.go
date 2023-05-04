/*
	Checks cryptographic algorithms correctness and execution time
*/

package singleFunctionTesting

import (
	"bavovnacoin/byteArr"
	"bavovnacoin/ecdsa"
	"bavovnacoin/hashing"
	"crypto/sha1"
	"fmt"
	"strings"
	"time"
)

type CryptoTest struct {
	testHashData []string

	Sha1TestAmmount  uint
	EcdsaTestAmmount uint
}

func (ct *CryptoTest) genText(isLong bool) string {
	if isLong {
		return strings.Repeat(fmt.Sprint(time.Now()), 15)
	}
	return fmt.Sprint(time.Now())
}

func (ct *CryptoTest) sha1Test() {
	var meanDuration [2]time.Duration
	var isPassed bool = true

	for i := uint(0); i < ct.Sha1TestAmmount; i++ {
		text := ct.genText(i%2 == 1)
		var customAlgo byteArr.ByteArr

		st := time.Now()
		customAlgo_ := hashing.SHA1(text)
		meanDuration[i%2] += time.Since(st) / time.Microsecond

		customAlgo.SetFromHexString(customAlgo_, 20)

		var pkgAlgo byteArr.ByteArr
		pkgAlgo_ := sha1.Sum([]byte(text))
		copy(pkgAlgo.ByteArr[:], pkgAlgo_[:])

		if !customAlgo.IsEqual(pkgAlgo) {
			isPassed = false
			println("Error for message", text)
		}
	}

	println("Mean time for short text (ns):", int64(meanDuration[0])/int64(ct.Sha1TestAmmount)/2)
	div := int64(ct.Sha1TestAmmount) / 2
	if ct.Sha1TestAmmount%2 == 1 {
		div = int64(ct.Sha1TestAmmount-1) / 2
	}
	println("Mean time for long text (ns):", int64(meanDuration[1])/div)
	if isPassed {
		println("Values are correct")
		println("Test passed")
	}
}

func (ct *CryptoTest) ecdsaTest() {
	var meanSignTime time.Duration
	var meanVerifTime time.Duration
	var isPassed bool = true

	ecdsa.InitValues()
	for i := uint(0); i < ct.Sha1TestAmmount; i++ {
		kp := ecdsa.GenKeyPair()
		hashText := hashing.SHA1(ct.genText(false))

		st := time.Now()
		s := ecdsa.Sign(hashText, kp.PrivKey)
		meanSignTime += time.Since(st) / time.Microsecond

		st = time.Now()
		isValid := ecdsa.Verify(kp.PublKey, s, hashText)
		meanVerifTime += time.Since(st) / time.Microsecond

		if !isValid {
			isPassed = false
			println("Error for hash mes", hashText)
		}
	}

	println("Mean time for signing (ns):", int64(meanSignTime)/int64(ct.EcdsaTestAmmount))
	println("Mean time for verifying (ns):", int64(meanVerifTime)/int64(ct.EcdsaTestAmmount))
	if isPassed {
		println("Values are correct")
		println("Test passed")
	}
}

func (ct *CryptoTest) Launch(algoName string) {
	if ct.Sha1TestAmmount <= 0 {
		ct.Sha1TestAmmount = 6
	}
	if ct.EcdsaTestAmmount <= 0 {
		ct.EcdsaTestAmmount = 6
	}

	if algoName == "SHA1" {
		ct.sha1Test()
	} else if algoName == "ECDSA" {
		ct.ecdsaTest()
	} else {
		println("Such an algorithm is not found")
	}
}
