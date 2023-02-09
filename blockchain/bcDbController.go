package blockchain

import (
	"bavovnacoin/dbController"
	"fmt"
)

func SetBcHeight(hgt uint64) bool {
	byteVal, isConv := dbController.ToByteArr(hgt)
	if !isConv {
		return false
	}
	return Database.SetValue("bcLength", byteVal)
}

func IncrBcHeight() {
	BcLength++
	SetBcHeight(BcLength)
}

func GetBcHeight() (uint64, bool) {
	value, isGotten := Database.GetValue("bcLength")
	if !isGotten {
		return 0, false
	}

	var len uint64
	isConv := dbController.FromByteArr(value, &len)
	if !isConv {
		return 0, false
	}
	return len, true
}

func WriteBlock(height uint64, block Block) bool {
	byteVal, isConv := dbController.ToByteArr(block)
	if !isConv {
		return false
	}
	println("Writing block", "bc"+fmt.Sprint(height))
	return Database.SetValue("bc"+fmt.Sprint(height), byteVal)
}

func GetBlock(height uint64) (Block, bool) {
	var block Block
	value, isGotten := Database.GetValue("bc" + fmt.Sprint(height))
	if !isGotten {
		return block, false
	}

	isConv := dbController.FromByteArr(value, &block)
	if !isConv {
		return block, false
	}
	return block, true
}
