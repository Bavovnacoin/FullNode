package blockchain

import (
	"bavovnacoin/hashing"
	"bavovnacoin/node_controller/command_executor"
	"fmt"
	"log"
	"math/big"
	"runtime"
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

	var step uint64 = 1
	var start uint64 = (uint64(data.thrId) + data.thrCount*(step-1)) * data.iterPerStep
	var end uint64 = (uint64(data.thrId+1) + data.thrCount*(step-1)) * data.iterPerStep

	for ; allowParallelMining; step++ {
		var nonce uint64 = start

		for ; nonce < end; nonce++ {
			if !allowParallelMining {
				data.isFound = false
				ch <- data
			}

			data.block.Nonce = nonce
			hashNounce, _ := new(big.Int).SetString(hashing.SHA1(BlockToString(data.block)+fmt.Sprintf("%d", nonce)), 16)
			// println(fmt.Sprintf("%x", target))
			// println(fmt.Sprintf("%x", hashNounce))
			if target.Cmp(hashNounce) == 1 {
				data.isFound = true
				data.nonce = nonce
				ch <- data
				allowParallelMining = false
			}
			command_executor.PauseCommand()
		}

		start = (uint64(data.thrId) + data.thrCount*(step-1)) * data.iterPerStep
		end = (uint64(data.thrId+1) + data.thrCount*(step-1)) * data.iterPerStep
		if command_executor.ShowMiningStats {
			fmt.Printf("Mining stats. Thread [%d] is done. Now mining in range [%d - %d]\n", data.thrId, start, end)
		}
	}
}

func MineThreads(block Block, threadsCount uint64, allowPrint bool) uint64 {
	if allowPrint {
		log.Println("Mining started")
	}

	allowParallelMining = true
	var thrcount uint64
	if threadsCount+4 > uint64(runtime.NumCPU()) && uint64(runtime.NumCPU())-4 >= 1 {
		thrcount = uint64(runtime.NumCPU()) - 4
		if allowPrint {
			log.Println("Threads for mining are limited to " + fmt.Sprint(thrcount))
		}
	} else if uint64(runtime.NumCPU())-4 < 1 {
		thrcount = 1
		if allowPrint {
			log.Println("Threads for mining are limited to 1")
		}
	} else {
		thrcount = threadsCount
	}

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
	if allowPrint {
		log.Println("Mining done")
	}
	return foundNounce
}
