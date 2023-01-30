package command_executor

import (
	"bavovnacoin/account"
	"time"
)

var Node_working bool = true
var Pause bool = false
var Network_accounts []account.Account

func PauseCommand() {
	for Pause {
		time.Sleep(500 * time.Millisecond)
	}
}
