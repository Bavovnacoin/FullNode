package blockchain

import (
	"bavovnacoin/byteArr"
	"bavovnacoin/dbController"
	"fmt"
	"strconv"

	"github.com/syndtr/goleveldb/leveldb/util"
)

func SetBcHeight(height uint64, chainId uint64) bool {
	byteVal, isConv := byteArr.ToByteArr(height)
	if !isConv {
		return false
	}
	return dbController.DB.SetValue("bcLength"+fmt.Sprintf("%d:", chainId), byteVal)
}

func IncrBcHeight(chainId uint64) {
	BcLength++
	SetBcHeight(BcLength, chainId)
}

func GetBcHeight(chainId uint64) (uint64, bool) {
	value, isGotten := dbController.DB.GetValue("bcLength" + fmt.Sprintf("%d:", chainId))
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

func getAllLastBlocks() ([]Block, []uint64, []uint64) {
	var heights []uint64
	var chainIds []uint64
	var height uint64

	iter := dbController.DB.Db.NewIterator(util.BytesPrefix([]byte("bcLength")), nil)
	for iter.Next() {
		byteArr.FromByteArr(iter.Value(), &height)
		heights = append(heights, height)
		chainIds = append(chainIds, getChainIdFromKey(string(iter.Key())))
	}
	iter.Release()

	var blocks []Block
	for i := 0; i < len(heights); i++ {
		block, _ := GetBlock(heights[i]-1, chainIds[i])
		blocks = append(blocks, block)
	}
	return blocks, chainIds, heights
}

func WriteBlock(height uint64, chainId uint64, block Block) bool {
	byteVal, isConv := byteArr.ToByteArr(block)
	if !isConv {
		return false
	}
	return dbController.DB.SetValue("bc"+fmt.Sprint(height)+":"+fmt.Sprint(chainId), byteVal)
}

func RemBlock(height uint64, chainId uint64) bool {
	return dbController.DB.Db.Delete([]byte("bc"+fmt.Sprint(height)+":"+fmt.Sprint(chainId)), nil) == nil
}

func GetBlock(height uint64, chainId uint64) (Block, bool) {
	var block Block

	value, isGotten := dbController.DB.GetValue("bc" + fmt.Sprint(height) + ":" + fmt.Sprint(chainId))
	if !isGotten {
		return block, false
	}

	isConv := byteArr.FromByteArr(value, &block)
	if !isConv {
		return block, false
	}
	return block, true
}

type BlockChainId struct {
	Block   Block
	ChainId uint64
}

func getChainIdFromKey(key string) uint64 {
	for i := len(key) - 1; i >= 0; i-- {
		if key[i] == ':' {
			num, _ := strconv.ParseUint(key[i+1:], 10, 64)
			return num
		}
	}
	return 0
}

func GetBlocksOnHeight(height uint64) ([]BlockChainId, bool) {
	var blockArr []BlockChainId
	var block_id BlockChainId

	iter := dbController.DB.Db.NewIterator(util.BytesPrefix([]byte("bc"+fmt.Sprint(height)+":")), nil)
	for iter.Next() {
		byteArr.FromByteArr(iter.Value(), &block_id.Block)
		block_id.ChainId = getChainIdFromKey(string(iter.Key()))
		blockArr = append(blockArr, block_id)
	}
	iter.Release()

	if iter.Error() != nil {
		return blockArr, false
	}

	return blockArr, true
}

func SetBlockForkHeight(height uint64, chainId uint64) bool {
	byteVal, isConv := byteArr.ToByteArr(height)
	if !isConv {
		return false
	}
	return dbController.DB.SetValue("forkHeight"+fmt.Sprintf("%d:", chainId), byteVal)
}

func GetBlockForkHeight(chainId uint64) (uint64, bool) {
	value, isGotten := dbController.DB.GetValue("forkHeight" + fmt.Sprintf("%d:", chainId))
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
