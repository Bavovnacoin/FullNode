package main

import "bavovnacoin/testing/singleFunctionTesting"

func main() {
	var test1 singleFunctionTesting.BlockchainVerifTest
	test1.SetTestValues(1, 1, 3)
	test1.BlockchainVerefication()
}
