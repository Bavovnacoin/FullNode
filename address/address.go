package address

import (
	"fmt"
	"strconv"
	"strings"
)

type Address struct {
	byteArr [20]byte
}

func addZerosAtBeginning(val string, totalLen int) string {
	return strings.Repeat("0", totalLen-len(val)) + val
}

func (addr *Address) SetFromHexString(hexVal string) bool {
	hexVal = addZerosAtBeginning(hexVal, 40)
	for i := 0; i < len(hexVal); i += 2 {
		numVal, err := strconv.ParseUint(hexVal[i:i+2], 16, 64)
		if err != nil {
			return false
		}
		addr.byteArr[i/2] = byte(numVal)
	}

	return true
}

func (addr Address) ToHexString() string {
	var res string
	for i := 0; i < len(addr.byteArr); i++ {
		res += addZerosAtBeginning(fmt.Sprintf("%x", addr.byteArr[i]), 2)
	}
	return res
}

func (addr Address) IsEqual(newAddr Address) bool {
	for i := 0; i < len(newAddr.byteArr); i++ {
		if addr.byteArr[i] != newAddr.byteArr[i] {
			return false
		}
	}
	return true
}
