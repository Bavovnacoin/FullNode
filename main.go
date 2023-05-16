package main

import (
	"bavovnacoin/testing/singleFunctionTesting"
	"fmt"
)

func ByteToString(arr []byte) string {
	str := ""
	for i := 0; i < len(arr); i++ {
		str += fmt.Sprint(arr[i])
	}
	return str
}

func main() {
	// node.Launch()
	var ct singleFunctionTesting.CommunicationTest
	ct.Launch()
}
