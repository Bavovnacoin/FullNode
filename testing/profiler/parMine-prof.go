package profiler

import (
	"bavovnacoin/testing/singleFunctionTesting"
	"net/http"
)

func LaunchServer(handler func(w http.ResponseWriter, req *http.Request)) {
	http.HandleFunc("/", handler)
	http.ListenAndServe("localhost:8080", nil)
}

func HandleParMining(w http.ResponseWriter, req *http.Request) {
	var ecdsa singleFunctionTesting.ParallelMiningTest
	ecdsa.Launch(5)
}
