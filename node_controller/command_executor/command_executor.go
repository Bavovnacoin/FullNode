package command_executor

import (
	"time"
)

var Node_working bool = true
var ShowMiningStats = false
var Pause bool = false

func PauseCommand() {
	for Pause {
		time.Sleep(500 * time.Millisecond)
	}
}
