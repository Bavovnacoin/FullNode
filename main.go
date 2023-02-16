package main

import "bavovnacoin/testing/singleFunctionTesting"

func isHex(s string) bool {
	arr := "0123456789abcdef"
	counter := 0
	for i := 0; i < len(s); i++ {
		for j := 0; j < len(arr); j++ {
			if s[i] == arr[j] {
				counter++
				break
			}
		}
	}
	if counter != len(s) {
		return false
	}
	return true
}

func main() {
	singleFunctionTesting.TransactionsVerefication()
}
