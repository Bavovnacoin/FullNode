package node_settings

import (
	"bavovnacoin/node_controller/command_executor"
	"fmt"
	"strconv"
)

var settings_menu bool

func fieldEnterForm(headerMes string, settings *NodeSettings, execFunc func(inp string, settings *NodeSettings) bool) {
	var errStr string
	var input string
	for true {
		command_executor.ComContr.ClearConsole()
		if errStr != "" {
			println(errStr)
			errStr = ""
		}

		fmt.Printf(headerMes)
		fmt.Scan(&input)
		if input == "back" {
			return
		} else {
			if execFunc(input, settings) {
				return
			} else {
				errStr = "Wrong input. Try again"
				continue
			}
		}
	}
}

func feePerByteValid(feeStr string, settings *NodeSettings) bool {
	fee, err := strconv.ParseUint(feeStr, 10, 64)
	if err == nil && fee >= 0 {
		settings.TxMinFee = fee
		settings.WriteSettings()
		return true
	}
	return false
}

func threadsSetValid(threadsStr string, settings *NodeSettings) bool {
	threads, err := strconv.ParseUint(threadsStr, 10, 32)
	if err == nil && settings.SetMiningThreads(uint(threads)) {
		settings.WriteSettings()
		return true
	}
	return false
}

func nodeNamesVariants(names []string, currNodeTypeID uint) string {
	res := ""
	outInd := 1
	for i := 0; i < len(names); i++ {
		if i != int(currNodeTypeID) {
			res += fmt.Sprintf("%d. %s\n", outInd, names[i])
			outInd++
		}
	}
	return res
}

func nodeTypeSetValid(nodeType string, settings *NodeSettings) bool {
	tp, err := strconv.ParseUint(nodeType, 10, 32)
	if err == nil && tp >= 0 && tp < uint64(len(settings.NodeTypesNames)) {
		if tp <= uint64(settings.NodeType) {
			settings.NodeType = uint(tp) - 1
		} else {
			settings.NodeType = uint(tp)
		}
		settings.WriteSettings()
		return true
	}
	return false
}

func addressesVariants(addresses []string) string {
	if len(addresses) == 0 {
		return "There's no addresses\n"
	}

	res := ""
	i := 0
	for ; i < len(addresses); i++ {
		res += fmt.Sprintf("%d. %s\n", i+1, addresses[i])
	}

	return res
}

func addressesSetValid(input string, settings *NodeSettings) bool {
	variant, err := strconv.Atoi(input)
	if err == nil && variant >= 0 && variant-1 < len(settings.OtherNodesAddresses) {
		var remVar string
		for true {
			fmt.Printf("Are you shure that you want to delete %s address?\n", settings.OtherNodesAddresses[variant-1])
			println("Type \"Yes\" or \"No\"")
			fmt.Scan(&remVar)

			if remVar == "Yes" {
				settings.RemAddress(variant - 1)
				settings.WriteSettings()
				return true
			} else {
				return false
			}
		}
	} else {
		isAdded := settings.AddAddress(input)
		if isAdded {
			settings.WriteSettings()
			return true
		}
	}
	return false
}

func ChooseMenuVariant(variant string, settings *NodeSettings) {
	if variant == "1" {
		fieldEnterForm(fmt.Sprintf("Current fee is %d. Type another fee or \"back\" to back to the settings menu.\n",
			settings.TxMinFee), settings, feePerByteValid)
	} else if variant == "2" {
		fieldEnterForm(fmt.Sprintf("Current threads ammount is %s. Type another ammount (0 - for maximum) or \"back\" to back to the settings menu.\n",
			settings.ThreadsForMiningToString()), settings, threadsSetValid)
	} else if variant == "3" {
		fieldEnterForm(fmt.Sprintf("Current node type is %s. Type number for node type or \"back\" to back to the settings menu:\n%s\n",
			settings.NodeTypesNames[settings.NodeType], nodeNamesVariants(settings.NodeTypesNames, settings.NodeType)), settings, nodeTypeSetValid)
	} else if variant == "4" {
		fieldEnterForm("Current addresses are seen below. Choose a variant to delete or type a new address to add a new address. Type \"back\" to back to the settings menu.\n"+
			addressesVariants(settings.OtherNodesAddresses), settings, addressesSetValid)
	} else if variant == "5" {
		settings_menu = false
	}
}

func LaunchMenu(settings *NodeSettings) {
	settings_menu = true
	var variant string

	for settings_menu {
		command_executor.ComContr.ClearConsole()
		println("Press a right button to manage an option")
		fmt.Printf("1. Transaction min fee per byte (%d bavov-kos).\n", settings.TxMinFee)
		fmt.Printf("2. Threads for mining (%s).\n", settings.ThreadsForMiningToString())
		fmt.Printf("3. Node type (%s).\n", settings.NodeTypesNames[settings.NodeType])
		fmt.Printf("4. Other node addresses (Ammount: %d).\n", len(settings.OtherNodesAddresses))
		fmt.Printf("5. Back\n")

		fmt.Scan(&variant)
		ChooseMenuVariant(variant, settings)
	}
}
