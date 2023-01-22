package main

import (
	"bavovnacoin/hashing"
	"bavovnacoin/mining"
	"fmt"
	"math/big"
	"runtime"
	"time"
)

// TODO: Change message to block in the future
func mine(message string, bits string) int64 {
	st := time.Now()
	target := mining.BitsToTarget(bits)

	var nounce int64 = 0
	for ; true; nounce++ {
		hashNounce, _ := new(big.Int).SetString(hashing.SHA1(message+fmt.Sprintf("%d", nounce)), 16)
		if target.Cmp(hashNounce) == 1 {
			elaps := time.Since(st).Milliseconds()
			println(elaps)
			return nounce
		}
	}
	elaps := time.Since(st).Milliseconds()
	println(nounce)
	println(elaps)
	return -1
}

var allowParallelMining bool

type ParMineData struct {
	start, end uint64
	bits       string
	message    string
	nounce     uint64
	isFound    bool
}

// TODO: Change message to block in the future
func mineParTask(data ParMineData, ch chan ParMineData) {
	target := mining.BitsToTarget(data.bits)
	var nounce uint64 = data.start

	for ; nounce < data.end; nounce++ {
		if !allowParallelMining {
			data.isFound = false
			ch <- data
		}

		hashNounce, _ := new(big.Int).SetString(hashing.SHA1(data.message+fmt.Sprintf("%d", nounce)), 16)
		if target.Cmp(hashNounce) == 1 {
			data.isFound = true
			data.nounce = nounce
			println(fmt.Sprintf("%x", target))
			println(fmt.Sprintf("%x", hashNounce))
			ch <- data
			allowParallelMining = false
		}
	}

	data.isFound = false
	ch <- data
}

func parallel(message string, bits string) {
	st := time.Now()
	allowParallelMining = true
	thrcount := runtime.NumCPU() // Set minus for IO thread
	resChan := make(chan ParMineData, thrcount)
	//var wg sync.WaitGroup

	var step uint64 = 1
	for ; allowParallelMining; step++ {
		//wg.Add(thrcount)
		var iterPerStep uint64 = 100000 * step // TODO: change here???

		for i := 0; i < thrcount-1; i++ {
			thrData := ParMineData{start: uint64(i) * iterPerStep, end: uint64(i+1) * iterPerStep,
				bits: bits, message: message}
			go mineParTask(thrData, resChan)
		}
		thrData := ParMineData{start: uint64(thrcount-1) * iterPerStep, end: uint64(thrcount) * iterPerStep,
			bits: bits, message: message}
		go mineParTask(thrData, resChan)

		for i := 0; i < thrcount; i++ {
			data := <-resChan
			if data.isFound {
				println(data.nounce)
				//break
			}
		}
		println(fmt.Sprintf("Step %d - not found", step))
		//wg.Wait()
	}
	//wg.Done()
	//close(resChan)
	elaps := time.Since(st).Milliseconds()
	println(elaps)
}

func main() {
	bits, _ := new(big.Int).SetString("000000000ffff000000000000000000000000000", 16)
	// println("Common")
	// println(mine("Hello worldd", mining.TargetToBits(bits)))

	println("\nParallel")
	parallel("Hello worldd", mining.TargetToBits(bits))

}
