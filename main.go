package main

import loadtesting "bavovnacoin/testing/loadTesting"

func main() {
	//node.Launch()

	// var tvt singleFunctionTesting.TransVerifTest
	// tvt.TransactionsVerefication(10, 4)

	// var bvt singleFunctionTesting.BlockchainVerifTest
	// bvt.BlockchainVerefication(10, 3, 10)

	var lt loadtesting.LoadTest
	lt.StartLoadTest(20, 5)
}
