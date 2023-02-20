package byteArr

import (
	"fmt"
	"strconv"
	"strings"
)

type ByteArr struct {
	ByteArr []byte
}

func addZerosAtBeginning(val string, totalLen int) string {
	return strings.Repeat("0", totalLen-len(val)) + val
}

func (arr *ByteArr) SetFromHexString(hexVal string, length int) bool {
	arr.ByteArr = make([]byte, length)
	hexVal = addZerosAtBeginning(hexVal, length*2)
	for i := 0; i < len(hexVal); i += 2 {
		numVal, err := strconv.ParseUint(hexVal[i:i+2], 16, 64)
		if err != nil {
			return false
		}
		arr.ByteArr[i/2] = byte(numVal)
	}

	return true
}

func (byteArr ByteArr) ToHexString() string {
	var res string
	for i := 0; i < len(byteArr.ByteArr); i++ {
		res += addZerosAtBeginning(fmt.Sprintf("%x", byteArr.ByteArr[i]), 2)
	}
	return res
}

func (arr ByteArr) IsEqual(newByteArr ByteArr) bool {
	for i := 0; i < len(newByteArr.ByteArr); i++ {
		if arr.ByteArr[i] != newByteArr.ByteArr[i] {
			return false
		}
	}
	return true
}
