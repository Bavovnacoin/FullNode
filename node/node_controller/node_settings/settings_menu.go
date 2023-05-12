package node_settings

import (
	"bavovnacoin/byteArr"
	"bavovnacoin/cryption"
	"bavovnacoin/ecdsa"
	"bavovnacoin/hashing"
	"bavovnacoin/node/node_controller/command_executor"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var settings_menu bool

func FieldEnterForm(headerMes string, settings *NodeSettings, execFunc func(inp string, settings *NodeSettings) bool) bool {
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
			return false
		} else {
			if execFunc(input, settings) {
				return true
			} else {
				errStr = "Wrong input. Try again"
				continue
			}
		}
	}
	return false
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

func addressesVariants(addresses [][]string) string {
	if len(addresses) == 0 {
		return "There's no addresses\n"
	}

	res := ""
	i := 0
	for ; i < len(addresses); i++ {
		res += fmt.Sprintf("%d. %s %s\n", i+1, addresses[i][0], addresses[i][1])
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
		addresses := strings.Split(input, "_")
		isAdded := settings.AddAddress(addresses)
		if isAdded {
			settings.WriteSettings()
			return true
		}
	}
	return false
}

func myAddrSetValid(address string, settings *NodeSettings) bool {
	if address != settings.MyAddress {
		settings.MyAddress = address
		settings.WriteSettings()
		return true
	}
	return false
}

func rpcAddrSetValid(address string, settings *NodeSettings) bool {
	if address != settings.RPCip {
		settings.RPCip = address
		settings.WriteSettings()
		return true
	}
	return false
}

func rewAddrSetValid(rewAddr string, settings *NodeSettings) bool {
	if rewAddr != settings.RewardAddress && settings.IsRewAddrWalid(rewAddr) {
		settings.RewardAddress = rewAddr
		settings.WriteSettings()
		return true
	}
	return false
}

const (
	PassSpecSign = " !”#$%%&()*+,./:;<=>?@[/^_`{}|~"
)

func IsPassValid(pass string) bool {
	// Regex for checking at least one lower, one upper, one num and spec sign
	reg := fmt.Sprintf("^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*[%s])[A-Za-z\\d%s]$", PassSpecSign, PassSpecSign)
	isPassMatch, _ := regexp.MatchString(reg, pass)
	if (len(pass) < 8 || len(pass) > 20) && !isPassMatch {
		return false
	}
	return true
}

func changePassword(input string, settings *NodeSettings) bool {
	var hash byteArr.ByteArr
	hash.SetFromHexString(hashing.SHA1(input), 20)

	if input == "NEW" {
		var password string
		for true {
			println("Enter a new password")
			fmt.Scan(&password)

			if IsPassValid(password) {
				break
			} else {
				println("The password is incorrect")
			}
			if password == "back" {
				return false
			}
		}

		var password_conf string
		for true {
			println("Confirm changes by typing in your new password")
			fmt.Scan(&password_conf)

			if password == password_conf {
				break
			}
			if password_conf == "back" {
				return false
			}
		}

		Settings.SetPrivKey(password)
		settings.HashPass.SetFromHexString(hashing.SHA1(password), 20)
		settings.WriteSettings()

	} else if settings.HashPass.IsEqual(hash) {
		var password string
		for true {
			println("Enter a new password")
			fmt.Scan(&password)

			if IsPassValid(password) {
				break
			} else {
				println("The password is incorrect")
			}
			if password == "back" {
				return false
			}
		}

		ecdsa.InitValues()
		prKey := cryption.AES_decrypt(string(Settings.PrivKey), input)
		Settings.PrivKey = []byte(cryption.AES_encrypt(prKey, password))

		var newPassHash byteArr.ByteArr
		newPassHash.SetFromHexString(hashing.SHA1(password), 20)
		Settings.HashPass = newPassHash
		Settings.WriteSettings()
		return true
	}

	return false
}

func ChooseMenuVariant(variant string, settings *NodeSettings) {
	if variant == "1" {
		FieldEnterForm(fmt.Sprintf("Current fee is %d. Type another fee or \"back\" to back to the settings menu.\n",
			settings.TxMinFee), settings, feePerByteValid)
	} else if variant == "2" {
		FieldEnterForm(fmt.Sprintf("Current threads ammount is %s. Type another ammount up to %d (0 - for maximum) or \"back\" to back to the settings menu.\n",
			settings.ThreadsForMiningToString(), settings.GetMaxThreadsAmmount()), settings, threadsSetValid)
	} else if variant == "3" {
		FieldEnterForm(fmt.Sprintf("Current node type is %s. Type number for node type or \"back\" to back to the settings menu:\n%s\n",
			settings.NodeTypesNames[settings.NodeType], nodeNamesVariants(settings.NodeTypesNames, settings.NodeType)), settings, nodeTypeSetValid)
	} else if variant == "4" {
		FieldEnterForm("Current addresses are seen below. Choose a variant to delete or type a new address to add a new address. Type \"back\" to back to the settings menu.\n"+
			addressesVariants(settings.OtherNodesAddresses), settings, addressesSetValid)
	} else if variant == "5" {
		FieldEnterForm(fmt.Sprintf("Current address is %s. Enter a new one or type \"back\" to back to the settings menu.\n", settings.MyAddress),
			settings, myAddrSetValid)
	} else if variant == "6" {
		FieldEnterForm(fmt.Sprintf("Current reward address is %s. Enter a new one or type \"back\" to back to the settings menu.\n", settings.GetRewAddress()),
			settings, rewAddrSetValid)
	} else if variant == "7" {
		FieldEnterForm(fmt.Sprintf("Current rpc address is %s. Enter a new one or type \"back\" to back to the settings menu.\n", settings.GetRpcAddr()),
			settings, rpcAddrSetValid)
	} else if variant == "8" {
		FieldEnterForm(fmt.Sprintf("Please, enter your old password or type \"NEW\" for a new password. WARNING: setting up a new password will cause node`s address change\n"),
			settings, changePassword)
	} else if variant == "0" {
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
		fmt.Printf("5. Change my address (Current: %s).\n", settings.MyAddress)
		fmt.Printf("6. Change reward address (Current: %s).\n", settings.GetRewAddress())
		fmt.Printf("7. Change RPC address (Current: %s).\n", settings.GetRpcAddr())
		fmt.Printf("8. Change the password.\n")

		fmt.Printf("0. Back\n")

		fmt.Scan(&variant)
		ChooseMenuVariant(variant, settings)
	}
}
