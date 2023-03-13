package main

import (
	"bavovnacoin/node"
	loadtesting "bavovnacoin/testing/loadTesting"
)

func main() {
	node.Launch()
	var lt loadtesting.LoadTest
	lt.StartLoadTest(22, 1, 5)
}
