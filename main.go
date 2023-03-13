package main

import loadtesting "bavovnacoin/testing/loadTesting"

func main() {
	//node.Launch()
	var lt loadtesting.LoadTest
	lt.StartLoadTest(50, 1, 10)
}
