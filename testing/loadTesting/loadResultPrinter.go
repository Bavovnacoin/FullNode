package loadtesting

import (
	"fmt"
	"time"
)

func printfunRes(funExecTime []time.Duration, funName string, timeSpecif uint64) {
	var funTimeChunks [3]float64
	ind := 0
	transStep := len(funExecTime) / 3
	fmt.Print(funName)
	var maxTime float64 = 0.000000000000000000000000000000001
	var minTime float64 = 99999999999999999999

	for i := 0; i < len(funExecTime); i++ {
		currTime := float64(funExecTime[i]) / float64(timeSpecif)
		if currTime > maxTime {
			maxTime = currTime
		}

		if currTime < minTime {
			minTime = currTime
		}
		funTimeChunks[ind] += currTime
		if i == (ind+1)*transStep && ind < 2 {
			fmt.Printf("%d/3 - %.4f", ind+1, funTimeChunks[ind]/float64(transStep))
			if ind != 2 {
				fmt.Print(", ")
			}
			ind++
		}
	}

	fmt.Printf("3/3 - %.4f. Max: %.4f, min: %.4f\n", funTimeChunks[2]/float64(transStep), maxTime, minTime)
}

func (lt *LoadTest) printResults() {
	println("Load test results.")
	println("RPC exec mean time results (ms):")
	printfunRes(lt.rpcExecTimeUtxoByAddr, "1. Get utxo by address: ", uint64(time.Millisecond))
	printfunRes(lt.rpcExecTimeisAddrExist, "2. Is address exists: ", uint64(time.Millisecond))
	println()
	println("New tx verification mean time results (s):")
	printfunRes(lt.txVerifTime, "Tx verification: ", uint64(time.Second))
}
