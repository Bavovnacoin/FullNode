package command_executor

import (
	"os"
	"os/exec"
	"time"
)

var ComContr CommandController

type CommandController struct {
	FullNodeWorking bool
	OpSys           string
}

// var FullNodeWorking bool = true
var ShowMiningStats = false
var Pause bool = false

func (cc *CommandController) ClearConsole() {
	if cc.OpSys == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else if cc.OpSys == "linux" || cc.OpSys == "darwin" {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		for i := 0; i < 3; i++ {
			println()
		}
	}
}

func (cc *CommandController) getLineSeparator() string {
	if cc.OpSys == "windows" {
		return "\r\n"
	} else if cc.OpSys == "darwin" {
		return "\r"
	} else {
		return "\n"
	}
}

func PauseCommand() {
	for Pause {
		time.Sleep(500 * time.Millisecond)
	}
}
