package blockchain

var startEmit uint64 = 7000000000

func GetCoinsForEmition() uint64 {
	return uint64(startEmit / ((uint64(BcLength+1)/2105280)/2 + 1))
}

func CheckEmitedCoins(reward uint64, blockchainLen int) bool {
	teorReward := uint64(startEmit / ((uint64(blockchainLen)/2105280)/2 + 1))
	return teorReward == reward
}
