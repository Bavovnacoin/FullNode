package node_controller

import (
	"bavovnacoin/account"
	"bavovnacoin/blockchain"
	"bavovnacoin/byteArr"
	"bavovnacoin/hashing"
	"bavovnacoin/node_controller/command_executor"
	"bavovnacoin/transaction"
	"bavovnacoin/txo"
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
)

var text string
var reader *bufio.Reader

func getLineSeparator() string {
	opSys := runtime.GOOS
	if opSys == "windows" {
		return "\r\n"
	} else if opSys == "darwin" {
		return "\r"
	} else {
		return "\n"
	}
}

func CommandHandler() {
	reader = bufio.NewReader(os.Stdin)
	for command_executor.Node_working {
		text, _ = reader.ReadString('\n')
		text = strings.Trim(text, getLineSeparator())
		if text == "stop" {
			command_executor.Node_working = false
			log.Println("Node has been stopped.")
		} else if text == "pause" {
			command_executor.Pause = !command_executor.Pause
			if command_executor.Pause {
				log.Println("Pause has been enabled.")
			} else {
				log.Println("Pause has been disabled.")
			}
		} else if text == "help" {
			helpPrinter()
		} else if text == "getmemp" {
			ifPauseFunction(blockchain.PrintMempool, "getmemp")
		} else if len(text) >= 9 && text[:9] == "getmemptx" {
			ifPauseFunction(mempoolTxPrinter, "getmemptx")
		} else if len(text) >= 5 && text[:5] == "getbc" {
			ifPauseFunction(bcPrinter, "getbc")
		} else if len(text) >= 10 && text[:10] == "getblocktx" {
			ifPauseFunction(blockTxPrinter, "getblocktx")
		} else if len(text) >= 10 && text[:10] == "getaccaddr" {
			ifPauseFunction(accAddressesPrinter, "getaccaddr")
		} else if text == "gettxo" {
			ifPauseFunction(utxoPrinter, "gettxo")
		} else if text == "maketx" {
			ifPauseFunction(makeTransaction, "maketx")
		} else if text == "showminingstats" {
			miningStatsPrinter()
		} else if text != "" {
			fmt.Printf("Command '%s' is unknown\n", text)
		}
	}
}

func helpPrinter() {
	log.Println("List of all commands:")
	println("stop - stop node")
	println("pause - pause node")
	println("help - get list of all commands")
	println("getmemp - show mempool transactions")
	println("getmemptx [id] - show specific mempool transaction")
	println("getbc [start_height] [end_height] - show titles of blockchain blocks from a defined range")
	println("getblocktx [block_height] [block_tx_id] - show specific transaction from a defined block")
	println("getaccaddr [acc_id] - show addresses and balances of a specific account")
	println("gettxo - show all outputs")
	println("maketx - create new transaction and send to mempool")
	println("showminingstats - show statistics of current mining process")
}

func ifPauseFunction(function func(), command string) {
	if command_executor.Pause {
		function()
	} else {
		fmt.Printf("Can't use '%s' command without pause", command)
	}
}

func mempoolTxPrinter() {
	command := strings.Split(text, " ")
	if len(command) == 1 {
		log.Println("Error. You typed in command without parameters")
		return
	}

	for i := 1; i < len(command); i++ {
		ind, err := strconv.ParseInt(command[i], 10, 64)
		if err == nil {
			if ind >= int64(len(blockchain.Mempool)) || ind < 0 {
				log.Println("Error. You typed wrong tx index. Must be between 0 and", len(blockchain.Mempool)-1)
				break
			}
			log.Println("Mempool transaction:")
			transaction.PrintTransaction(blockchain.Mempool[ind])
			break
		} else {
			log.Println("Error. Expected numerical type as a parameter")
		}
	}
}

func bcPrinter() {
	command := strings.Split(text, " ")
	if len(command) < 3 {
		log.Println("Error. You typed in command without parameters")
		return
	}
	var commandValues []int64

	for i := 1; i < len(command); i++ {
		ind, err := strconv.ParseInt(command[i], 10, 64)
		commandValues = append(commandValues, ind)

		if err == nil {
			if len(commandValues) == 2 {
				if commandValues[0] < 0 || commandValues[1] < 0 ||
					commandValues[0] > commandValues[1] || commandValues[0] > int64(blockchain.BcLength) {
					log.Println("Error. You typped in wrong index.")
					return
				}

				if commandValues[1] > int64(blockchain.BcLength) {
					commandValues[1] = int64(blockchain.BcLength)
				}

				log.Println("Blockchain from", commandValues[0], "to", commandValues[1])
				for i := commandValues[0]; i < commandValues[1]; i++ {
					block, _ := blockchain.GetBlock(uint64(i))
					blockchain.PrintBlockTitle(block, int(i))
					println()
				}
			}
		} else {
			log.Println("Error. Expected numerical type as a parameter")
		}
	}
}

func blockTxPrinter() {
	command := strings.Split(text, " ")
	if len(command) < 3 {
		log.Println("Error. You typed in command without parameters")
		return
	}
	var commandValues []int64

	for i := 1; i < len(command); i++ {
		ind, err := strconv.ParseInt(command[i], 10, 64)
		commandValues = append(commandValues, ind)

		if err == nil {
			if len(commandValues) == 2 {
				block, _ := blockchain.GetBlock(uint64(commandValues[0]))
				if commandValues[0] < 0 || commandValues[1] < 0 ||
					commandValues[0] > int64(blockchain.BcLength) ||
					commandValues[1] >= int64(len(block.Transactions)) {
					log.Println("Error. You typped in wrong index.")
					return
				}
				if commandValues[1] > int64(blockchain.BcLength) {
					commandValues[1] = int64(blockchain.BcLength)
				}
				log.Println("Block", commandValues[0], "- transaction", commandValues[1])
				transaction.PrintTransaction(block.Transactions[commandValues[1]])
			}
		} else {
			log.Println("Error. Expected numerical type as a parameter")
		}
	}
}

func accAddressesPrinter() {
	command := strings.Split(text, " ")
	if len(command) == 1 {
		log.Println("Error. You typed in command without parameters")
		return
	}

	for i := 1; i < len(command); i++ {
		ind, err := strconv.ParseInt(command[i], 10, 64)
		if err == nil {
			if ind >= int64(len(account.Wallet)) || ind < 0 {
				log.Println("Error. You typed wrong tx index. Must be between 0 and", len(account.Wallet)-1)
				break
			}

			log.Println("Address of account index", ind)
			acc := account.Wallet[ind]
			var value uint64
			for i := 0; i < len(acc.KeyPairList); i++ {
				var accAddress byteArr.ByteArr
				accAddrStr := hashing.SHA1(acc.KeyPairList[i].PublKey)
				accAddress.SetFromHexString(accAddrStr, 20)
				bal := account.GetBalByAddress(accAddress)
				value += bal
				fmt.Printf("[%d]. %s, balance: %d\n", i, accAddrStr, bal)
			}
			fmt.Printf("Total value: %d\n", value)
			break
		} else {
			log.Println("Error. Expected numerical type as a parameter")
		}
	}
}

func utxoPrinter() {
	txo.PrintCoinDatabase()
}

func makeTransaction() {
	log.Println("Transaction creation. To stop type stopcreation")

	var accId uint64
	println("Type in account id")
	for true {
		text, _ = reader.ReadString('\n')
		text = strings.Trim(text, getLineSeparator())
		if text == "stopcreation" {
			return
		}

		accIdInp, err := strconv.ParseUint(text, 10, 64)
		if err == nil {
			accId = accIdInp
			account.CurrAccount = account.Wallet[accId]
			break
		} else {
			println("Error. Expected numerical value.")
		}
	}

	var outAddr []byteArr.ByteArr
	var outValue []uint64
	println("Type in address and value to be sent separated by a space. Or continue by typing next.")

	for true {
		text, _ = reader.ReadString('\n')
		text = strings.Trim(text, getLineSeparator())
		inputArr := strings.Split(text, " ")
		if text == "next" {
			break
		}
		if text == "stopcreation" {
			return
		}
		var inpAddr byteArr.ByteArr
		inpAddr.SetFromHexString(inputArr[0], 20)
		for i := 1; i < len(inputArr); i++ {
			value, err := strconv.ParseUint(inputArr[i], 10, 64)
			if err == nil {
				outValue = append(outValue, value)
			}
		}
	}

	println("Type in tx fee per byte.")
	var fee int
	for true {
		text, _ = reader.ReadString('\n')
		text = strings.Trim(text, getLineSeparator())
		if text == "stopcreation" {
			return
		}

		feeInp, err := strconv.ParseInt(text, 10, 64)
		if err == nil {
			fee = int(feeInp)
			break
		} else {
			println("Error. Expected numerical value.")
		}
	}

	println("Type in locktime.")
	var locktime uint64
	for true {
		text, _ = reader.ReadString('\n')
		text = strings.Trim(text, getLineSeparator())
		if text == "stopcreation" {
			return
		}

		locktimeInp, err := strconv.ParseUint(text, 10, 64)
		if err == nil {
			locktime = locktimeInp
			break
		} else {
			println("Error. Expected numerical value.")
		}
	}

	tx, mes := transaction.CreateTransaction(fmt.Sprint(accId), outAddr, outValue, fee, uint(locktime))
	if mes == "" {
		transaction.PrintTransaction(tx)
		if blockchain.AddTxToMempool(tx, true) {
			log.Println("Created transaction is sent to the mempool")
		} else {
			log.Println("Transaction was not added to mempool")
		}
	} else {
		println(mes)
	}
}

func miningStatsPrinter() {
	command_executor.ShowMiningStats = !command_executor.ShowMiningStats
	if command_executor.ShowMiningStats {
		log.Println("Mining statistics is enabled")
	} else {
		log.Println("Mining statistics is disabled")
	}
}
