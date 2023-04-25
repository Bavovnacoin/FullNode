package node

import (
	"bavovnacoin/node/node_audithor"
	"bavovnacoin/node/node_controller/command_executor"
	"bavovnacoin/node/node_controller/node_settings"
	"bavovnacoin/node/node_validator"
	"fmt"
	"runtime"
	"strings"
)

var NodeLaunched bool

func funcChoser(variant string, isNodeLaunchAllowed bool) {
	if variant == "1" && isNodeLaunchAllowed {
		command_executor.ComContr.ClearConsole()
		if node_settings.Settings.NodeType == 0 {
			node_validator.LaunchValidatorNode()
		} else if node_settings.Settings.NodeType == 1 {
			node_audithor.LaunchAudithor()
		}

	} else if variant == "2" && isNodeLaunchAllowed || variant == "1" && !isNodeLaunchAllowed {
		node_settings.LaunchMenu(&node_settings.Settings)
	} else if variant == "3" && isNodeLaunchAllowed || variant == "2" && !isNodeLaunchAllowed {
		NodeLaunched = false
	}
}

func getNodeLaunchSettingsError() (string, string) {
	nodeType := ""
	var errMess []string
	if node_settings.Settings.RewardAddress == "" {
		errMess = append(errMess, "reward address (1-6)")
	}
	if node_settings.Settings.MyAddress == "" {
		errMess = append(errMess, "node address (1-5)")
	}

	if node_settings.Settings.NodeType == 0 {
		nodeType = "a validator"
	} else if node_settings.Settings.NodeType == 1 {
		nodeType = "an audithor"
	}
	return strings.Join(errMess, ", "), nodeType
}

func Launch() {
	NodeLaunched = true
	command_executor.ComContr.OpSys = runtime.GOOS
	node_settings.Settings.GetSettings()
	node_settings.Settings.InitSettingsValues()

	var variant string
	for NodeLaunched {
		command_executor.ComContr.ClearConsole()
		println("Choose a variant and press the right button")
		launchMesErr, nodeType := getNodeLaunchSettingsError()
		if launchMesErr != "" {
			fmt.Printf("Can't start a node. You need to manage: %s\n", launchMesErr)
		}

		var btn int = 1
		if launchMesErr == "" {
			fmt.Printf("%d. Launch %s node\n", btn, nodeType)
			btn++
		}
		fmt.Printf("%d. Manage settings\n", btn)
		fmt.Printf("%d. Exit\n", btn+1)

		fmt.Scan(&variant)
		funcChoser(variant, launchMesErr == "")
	}

	command_executor.ComContr.ClearConsole()
	println("Thank you for supporting Bavovnacoin network!")
}
