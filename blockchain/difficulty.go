package blockchain

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"
)

var BLOCK_DIFF_CHECK int = 3
var BLOCK_CREATION_SEC int = 60
var STARTBITS uint64 = 0xffff13 //2

func GetCurrBitsValue() uint64 {
	var bits uint64

	if len(Blockchain)%BLOCK_DIFF_CHECK == 0 && len(Blockchain) != 0 {
		bits = GenBits(Blockchain[len(Blockchain)-BLOCK_DIFF_CHECK].Time,
			Blockchain[len(Blockchain)-1].Time, Blockchain[len(Blockchain)-1].Bits)
	} else if len(Blockchain) != 0 {
		bits = Blockchain[len(Blockchain)-1].Bits
	} else {
		bits = STARTBITS
	}
	return bits
}

func GenBits(frstBlockTime time.Time, secBlockTime time.Time, bits uint64) uint64 {
	spentTimeSec := secBlockTime.Unix() - frstBlockTime.Unix()
	expextTimeSec := BLOCK_DIFF_CHECK * BLOCK_CREATION_SEC
	coef := float64(spentTimeSec) / float64(expextTimeSec)

	target := BitsToTarget(bits)
	target = target.Mul(target, big.NewInt(int64(coef*10000)))
	target = target.Div(target, big.NewInt(10000))
	targetStr := fmt.Sprintf("%x", target)

	if len(targetStr)%2 != 0 {
		targetStr = "0" + targetStr
	}

	targetNum, _ := new(big.Int).SetString(targetStr[4:6]+targetStr[:4]+strings.Repeat("0", len(targetStr)-6), 16)
	return TargetToBits(targetNum)
}

func addZerToLength(mes string, length int) string {
	return strings.Repeat("0", length-len(mes)) + mes
}

func BitsToTarget(bits uint64) *big.Int {
	bitsStr := addZerToLength(fmt.Sprintf("%x", bits), 8)
	println(bitsStr, "bits str")
	shift, _ := new(big.Int).SetString(bitsStr[6:], 16)
	shift.Sub(shift, big.NewInt(3))
	shift.Mul(shift, big.NewInt(8))
	powBase := big.NewInt(2)
	shift = powBase.Exp(powBase, shift, nil)
	println(fmt.Sprintf("%x - shift", shift))

	target, _ := new(big.Int).SetString(bitsStr[2:6]+bitsStr[:2], 16)
	target.Mul(target, shift)
	println(fmt.Sprintf("%x - mul", target))
	return target
}

func TargetToBits(target *big.Int) uint64 {
	targetStr := fmt.Sprintf("%x", target)
	var targetShift string
	var valStr string
	if len(targetShift)%2 != 0 {
		targetShift = targetStr[5:]
		valStr = targetStr[3:5] + "0" + targetStr[:3]
	} else {
		targetShift = targetStr[6:]
		valStr = targetStr[4:6] + targetStr[:4]
	}
	shiftVal, _ := big.NewInt(0).SetString("1"+targetShift, 16)
	shiftValFloat, _ := strconv.ParseFloat(shiftVal.String(), 40)

	shift := fmt.Sprintf("%x", (int(math.Log2(shiftValFloat)/8))+3)
	res, _ := strconv.ParseUint(valStr+shift, 16, 64)
	return res
}
