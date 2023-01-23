package blockchain

import (
	"bavovnacoin/hashing"
	"fmt"
	"math/big"
	"runtime"
)

var DIFF_CHECK_HOURS int = 24
var BLOCK_CREATION_SEC int = 60
var STARTBITS uint64 = 0xffff12

func MineCommon(block Block, bits uint64) uint64 {
	message := BlockToString(block)
	target := BitsToTarget(bits)

	var nounce uint64 = 0
	for ; true; nounce++ {
		block.Nonce = nounce
		hashNounce, _ := new(big.Int).SetString(hashing.SHA1(message+fmt.Sprintf("%d", nounce)), 16)
		if target.Cmp(hashNounce) == 1 {
			return nounce
		}
	}

	return 0
}

var allowParallelMining bool

type ParMineData struct {
	start, end uint64
	bits       uint64
	block      Block
	nounce     uint64
	isFound    bool
}

func mineParTask(data ParMineData, ch chan ParMineData) {
	target := BitsToTarget(data.bits)
	var nounce uint64 = data.start

	for ; nounce < data.end; nounce++ {
		if !allowParallelMining {
			data.isFound = false
			ch <- data
		}
		data.block.Nonce = nounce
		hashNounce, _ := new(big.Int).SetString(hashing.SHA1(BlockToString(data.block)+fmt.Sprintf("%d", nounce)), 16)
		if target.Cmp(hashNounce) == 1 {
			data.isFound = true
			data.nounce = nounce
			ch <- data
			allowParallelMining = false
		}
	}

	data.isFound = false
	ch <- data
}

func MineAllThreads(block Block, bits uint64) uint64 {
	allowParallelMining = true
	thrcount := uint64(runtime.NumCPU()) // Set minus for IO thread
	resChan := make(chan ParMineData, thrcount)
	var foundNounce uint64

	var step uint64 = 1
	for ; allowParallelMining; step++ {
		var iterPerStep uint64 = 10000

		var i uint64 = 0
		for ; i < thrcount-1; i++ {
			thrData := ParMineData{start: (uint64(i) + thrcount*(step-1)) * iterPerStep,
				end:  (uint64(i+1) + thrcount*(step-1)) * iterPerStep,
				bits: bits, block: block}
			go mineParTask(thrData, resChan)
		}
		thrData := ParMineData{start: (uint64(i) + thrcount*(step-1)) * iterPerStep,
			end:  (uint64(i+1) + thrcount*(step-1)) * iterPerStep,
			bits: bits, block: block}
		go mineParTask(thrData, resChan)

		i = 0
		for ; i < thrcount; i++ {
			data := <-resChan
			if data.isFound {
				foundNounce = data.nounce
			}
		}
	}

	return foundNounce
}
