package main

import (
	"bavovnacoin/testing/loadTesting/nodeLoadTest"
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
	var lt nodeLoadTest.LoadTest
	lt.Launch(100, 10)
}
