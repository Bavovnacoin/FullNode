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
	var execTime time.Duration
	var meanDuration [2]time.Duration
	var isPassed bool = true
	var passedAmmount uint

	for i := uint(0); i < ct.Sha1TestAmmount; i++ {
		text := ct.genText(i%2 == 1)
		var customAlgo byteArr.ByteArr

		st := time.Now()
		customAlgo_ := hashing.SHA1(text)

		execTime = time.Since(st) / time.Microsecond
		meanDuration[i%2] += execTime

		customAlgo.SetFromHexString(customAlgo_, 20)

		var pkgAlgo byteArr.ByteArr
		pkgAlgo_ := sha1.Sum([]byte(text))
		copy(pkgAlgo.ByteArr[:], pkgAlgo_[:])

		if customAlgo.IsEqual(pkgAlgo) {
			passedAmmount++
			textType := "Short text"
			if i%2 == 1 {
				textType = "Long text"
			}

			fmt.Printf("[%d]. Passed. %s. Time: %d mcs\n", i+1, textType, execTime)
		} else {
			isPassed = false
			println("Error for message", text)
		}
	}

	println("Mean time for short text (mcs):", int64(meanDuration[0])/int64(ct.Sha1TestAmmount)/2)
	div := int64(ct.Sha1TestAmmount) / 2
	if ct.Sha1TestAmmount%2 == 1 {
		div = int64(ct.Sha1TestAmmount-1) / 2
	}
	println("Mean time for long text (mcs):", int64(meanDuration[1])/div)

	if isPassed {
		fmt.Printf("Test passed (%d/%d)!\n", passedAmmount, ct.Sha1TestAmmount)
	} else {
		fmt.Printf("Test is not passed (%d/%d)!\n", passedAmmount, ct.Sha1TestAmmount)
	}
}

func (ct *CryptoTest) ecdsaTest() {
	var execSignTime time.Duration
	var execVerifTime time.Duration
	var meanSignTime time.Duration
	var meanVerifTime time.Duration
	var isPassed bool = true
	var passedAmmount uint

	ecdsa.InitValues()
	for i := uint(0); i < ct.Sha1TestAmmount; i++ {
		kp := ecdsa.GenKeyPair()
		hashText := hashing.SHA1(ct.genText(false))

		st := time.Now()
		s := ecdsa.Sign(hashText, kp.PrivKey)
		execSignTime = time.Since(st) / time.Microsecond
		meanSignTime += execSignTime

		st = time.Now()
		isValid := ecdsa.Verify(kp.PublKey, s, hashText)
		execVerifTime = time.Since(st) / time.Microsecond
		meanVerifTime += execVerifTime

		if isValid {
			passedAmmount++

			fmt.Printf("[%d]. Passed. Signing time: %d mcs, verification time: %d mcs.\n", i+1, execSignTime, execVerifTime)
		} else {
			isPassed = false
			println("Error for hash mes", hashText)
		}
	}

	println("Mean time for signing (mcs):", int64(meanSignTime)/int64(ct.EcdsaTestAmmount))
	println("Mean time for verifying (mcs):", int64(meanVerifTime)/int64(ct.EcdsaTestAmmount))
	if isPassed {
		fmt.Printf("Test passed (%d/%d)!\n", passedAmmount, ct.Sha1TestAmmount)
	} else {
		fmt.Printf("Test is not passed (%d/%d)!\n", passedAmmount, ct.Sha1TestAmmount)
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
