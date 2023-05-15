package node

import (
	"bavovnacoin/byteArr"
	"bavovnacoin/hashing"
	"bavovnacoin/node/node_audithor"
	"bavovnacoin/node/node_controller"
	"bavovnacoin/node/node_controller/command_executor"
	"bavovnacoin/node/node_settings"
	"bavovnacoin/node/node_validator"
	"fmt"
	"runtime"
	"strings"
)

var NodeLaunched bool

func funcChoser(variant string, isNodeLaunchAllowed bool) {
	if variant == "1" && isNodeLaunchAllowed && node_settings.Settings.HashPass.ByteArr != nil {
		command_executor.ComContr.ClearConsole()
		if node_settings.Settings.NodeType == 0 {
			node_validator.LaunchValidatorNode()
		} else if node_settings.Settings.NodeType == 1 {
			node_audithor.LaunchAudithor()
		}

	} else if variant == "2" && isNodeLaunchAllowed || variant == "1" && !isNodeLaunchAllowed {
		node_controller.LaunchSettingsMenu(&node_settings.Settings)
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
	if node_settings.Settings.HashPass.ByteArr == nil {
		errMess = append(errMess, "a password (TODO)")
	}

	if node_settings.Settings.NodeType == 0 {
		nodeType = "a validator"
	} else if node_settings.Settings.NodeType == 1 {
		nodeType = "an audithor"
	}
	return strings.Join(errMess, ", "), nodeType
}

func Authenticate() bool {
	var password string
	isPassInc := false

	for true {
		if isPassInc {
			println("The password is incorrect")
		}

		println("Please, enter your password to enter the system (or type \"exit\" to exit)")
		fmt.Scan(&password)

		var hashPass byteArr.ByteArr
		hashPass.SetFromHexString(hashing.SHA1(password), 20)

		if node_settings.Settings.HashPass.IsEqual(hashPass) {
			return true
		} else if password == "exit" {
			return false
		}
		isPassInc = true
	}
	return false
}

func passwordSetValid(password string, settings *node_settings.NodeSettings) bool {
	if !node_controller.IsPassValid(password) {
		return false
	}

	settings.HashPass.SetFromHexString(hashing.SHA1(password), 20)
	node_settings.Settings.SetPrivKey(password)
	settings.WriteSettings()
	return true
}

func Launch() {
	NodeLaunched = true
	command_executor.ComContr.OpSys = runtime.GOOS
	isNewFile := !node_settings.Settings.GetSettings()
	node_settings.Settings.InitSettingsValues()

	isPassSet := false
	if isNewFile {
		isPassSet = node_controller.FieldEnterForm(fmt.Sprintf("Please, enter a password.\nThe password length must be in range from 8 to 20 symbols. "+
			"It should contain upper and lower case letters, numbers and special signs (%s).\n", node_controller.PassSpecSign),
			&node_settings.Settings, passwordSetValid)
	}

	if !isPassSet && isNewFile {
		NodeLaunched = false
	}

	var variant string
	for NodeLaunched {
		command_executor.ComContr.ClearConsole()
		if !isPassSet {
			if !Authenticate() {
				break
			}
			command_executor.ComContr.ClearConsole()
			isPassSet = true
		}

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
