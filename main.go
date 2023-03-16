package main

import (
	"bavovnacoin/testing/singleFunctionTesting"
)

func main() {
	//node.Launch()

	// var tvt singleFunctionTesting.TransVerifTest
	// tvt.TransactionsVerefication(10, 4)

	// var bvt singleFunctionTesting.BlockchainVerifTest
	// bvt.BlockchainVerefication(10, 3, 10)

	var tvt singleFunctionTesting.TransVerifTime
	tvt.TransVerifTime()
}
