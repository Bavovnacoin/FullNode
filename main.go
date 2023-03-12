package main

import loadtesting "bavovnacoin/testing/loadTesting"

func main() {
	//node.Launch()
	var lt loadtesting.LoadTest
	lt.StartLoadTest(70, 1, 5)
}
