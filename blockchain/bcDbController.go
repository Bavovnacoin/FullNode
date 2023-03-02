package blockchain

import (
	"bavovnacoin/byteArr"
	"bavovnacoin/dbController"
	"fmt"
)

func SetBcHeight(hgt uint64) bool {
	byteVal, isConv := byteArr.ToByteArr(hgt)
	if !isConv {
		return false
	}
	return dbController.DB.SetValue("bcLength", byteVal)
}

func IncrBcHeight() {
	BcLength++
	SetBcHeight(BcLength)
}

func GetBcHeight() (uint64, bool) {
	value, isGotten := dbController.DB.GetValue("bcLength")
	if !isGotten {
		return 0, false
	}

	var len uint64
	isConv := byteArr.FromByteArr(value, &len)
	if !isConv {
		return 0, false
	}
	return len, true
}

func WriteBlock(height uint64, block Block) bool {
	byteVal, isConv := byteArr.ToByteArr(block)
	if !isConv {
		return false
	}
	return dbController.DB.SetValue("bc"+fmt.Sprint(height), byteVal)
}

func GetBlock(height uint64) (Block, bool) {
	var block Block
	value, isGotten := dbController.DB.GetValue("bc" + fmt.Sprint(height))
	if !isGotten {
		return block, false
	}

	isConv := byteArr.FromByteArr(value, &block)
	if !isConv {
		return block, false
	}
	return block, true
}
