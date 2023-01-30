package blockchain

import (
	"bavovnacoin/hashing"
	"bavovnacoin/node_controller/command_executor"
	"fmt"
	"math/big"
)

var allowParallelMining bool

type ParMineData struct {
	thrId, thrCount uint64
	iterPerStep     uint64
	bits            uint64
	block           Block
	nonce           uint64
	isFound         bool
}

func mineParTask(data ParMineData, ch chan ParMineData) {
	target := BitsToTarget(data.bits)
	//var nonce uint64 = data.start

	var step uint64 = 1
	for ; allowParallelMining; step++ {
		start := (uint64(data.thrId) + data.thrCount*(step-1)) * data.iterPerStep
		end := (uint64(data.thrId+1) + data.thrCount*(step-1)) * data.iterPerStep
		var nonce uint64 = start

		for ; nonce < end; nonce++ {
			if !allowParallelMining {
				data.isFound = false
				ch <- data
			}

			data.block.Nonce = nonce
			hashNounce, _ := new(big.Int).SetString(hashing.SHA1(BlockToString(data.block)+fmt.Sprintf("%d", nonce)), 16)
			if target.Cmp(hashNounce) == 1 {
				data.isFound = true
				data.nonce = nonce
				ch <- data
				allowParallelMining = false
			}
			command_executor.PauseCommand()
		}
	}
}

func MineThreads(block Block, threadsCount uint64) uint64 {
	fmt.Println("Mining started")
	allowParallelMining = true
	thrcount := threadsCount
	resChan := make(chan ParMineData, thrcount)
	var foundNounce uint64

	var iterPerStep uint64 = 10000

	var i uint64 = 0
	for ; i < thrcount-1; i++ {
		thrData := ParMineData{thrId: i, thrCount: thrcount, iterPerStep: iterPerStep,
			bits: block.Bits, block: block}
		go mineParTask(thrData, resChan)
	}
	thrData := ParMineData{thrId: i, thrCount: thrcount, iterPerStep: iterPerStep,
		bits: block.Bits, block: block}
	go mineParTask(thrData, resChan)

	i = 0
	for ; i < thrcount; i++ {
		data := <-resChan
		if data.isFound {
			foundNounce = data.nonce
		}
	}
	fmt.Println("Mining done")
	return foundNounce
}
