package mining

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"
)

// Difficulty changes every 24h
func GenBits(frstBlockTime time.Time, secBlockTime time.Time, bits string) {
	spentTimeSec := secBlockTime.Unix() - frstBlockTime.Unix()
	expextTimeSec := 24 * 60 * 60 // expected time for one day change
	coef := float64(expextTimeSec) / float64(spentTimeSec)
	target := BitsToTarget(bits)
	target = target.Mul(target, big.NewInt(int64(coef*100)))
	target = target.Div(target, big.NewInt(100))
	targetStr := fmt.Sprintf("%x", target)

	if len(targetStr)%2 != 0 {
		targetStr = "0" + targetStr
	}
	targetNum, _ := new(big.Int).SetString(targetStr[4:6]+targetStr[:4]+strings.Repeat("0", len(targetStr)-6), 16)
	println(TargetToBits(targetNum))
}

func BitsToTarget(bits string) *big.Int {
	shift, _ := new(big.Int).SetString(bits[6:], 16)
	shift.Sub(shift, big.NewInt(3))
	shift.Mul(shift, big.NewInt(8))
	powBase := big.NewInt(2)
	shift = powBase.Exp(powBase, shift, nil)

	target, _ := new(big.Int).SetString(bits[2:6]+bits[:2], 16)
	target.Mul(target, shift)
	return target
}

func TargetToBits(target *big.Int) string {
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
	shiftValFloat, _ := strconv.ParseFloat(shiftVal.String(), 64)

	shift := fmt.Sprintf("%x", (int(math.Log2(shiftValFloat)/8))+3)
	return valStr + shift
}
