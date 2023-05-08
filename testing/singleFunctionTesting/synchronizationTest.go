/*
	Tests how blocks are downloaded and uploaded (including altchains)
*/

package singleFunctionTesting

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/hashing"
)

type SynchronizationTest struct {
	SingleFunctionTesting

	mcBlocksAmmount int
	acBlocksAmmount int
	acAmmount       int
}

func (st SynchronizationTest) genAltchain(chainId uint64) {
	startHeight := uint64(st.mcBlocksAmmount) / 2
	prevBlock, _ := blockchain.GetBlock(startHeight, 0)

	newBlock := blockchain.Block{}
	for i := 0; i < st.acBlocksAmmount; i++ {
		newBlock, _ = st.CreateBlock(blockchain.GetBits(false),
			hashing.SHA1(blockchain.BlockHeaderToString(prevBlock)), prevBlock, false)
		newBlock.Version = 199

		blockchain.TryCameBlockToAdd(newBlock, startHeight+uint64(i)+1, []int64{})
		prevBlock = newBlock
		println()
	}
}

func (st SynchronizationTest) genMainchain() {
	newBlock := blockchain.Block{}
	for i := 0; i < st.mcBlocksAmmount; i++ {
		newBlock, _ = st.CreateBlock(blockchain.GetBits(false),
			hashing.SHA1(blockchain.BlockHeaderToString(blockchain.LastBlock)), blockchain.LastBlock, false)

		blockchain.AddBlockToBlockchain(newBlock, 0, true)
		blockchain.IncrBcHeight(0)
		blockchain.LastBlock = newBlock
	}
}

func (st SynchronizationTest) genChains() {

}

func (st SynchronizationTest) Launch() {
	st.mcBlocksAmmount = 10
	st.acBlocksAmmount = 2
	st.acAmmount = 1
	blockchain.STARTBITS = 0xffff14

	InitTestDb(true)
	st.genMainchain()
	st.genAltchain(1)

	// bl, res := blockchain.GetBlock(3, 0)
	// println(res)
	// blockchain.PrintBlockTitle(bl, 3)
}
