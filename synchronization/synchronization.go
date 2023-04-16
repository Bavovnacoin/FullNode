package synchronization

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/networking"
	"bavovnacoin/node_controller/node_settings"
	"log"
)

func StartSync(printLog bool, startBlock uint64) bool {
	var conn networking.Connection
	var addrInd int = -1

	var isEstablished bool
	isEstablished, addrInd = conn.EstablishAddresses(node_settings.Settings.OtherNodesAddresses, addrInd, "")
	if !isEstablished {
		return false
	}

	InitCheckpoints()
	var blockReqInd uint64 = startBlock
	checkpCorresp := false
	for true {
		blocks, currBcHeight, res := conn.RequestBlocks(blockReqInd)
		if res {
			if len(blocks) == 0 {
				return true
			}

			var blocksDownlSuccess bool = true
			for i := 0; i < len(blocks); i++ {
				checkpCorresp = checkForBlockCorrespondence(blockReqInd, blocks[i])
				if checkpCorresp {
					blockchain.AddBlockToBlockchain(blocks[i], 0) // TODO: make synchronization with altchains
					blockchain.IncrBcHeight(0)                    // TODO: make synchronization with altchains (height)
					blockReqInd++
				} else {
					if printLog {
						log.Printf("Address %s sent an incorrect block. Selecting next address\n",
							node_settings.Settings.OtherNodesAddresses[addrInd])
					}
					blocksDownlSuccess = false
					break
				}
			}
			if blocksDownlSuccess && printLog {
				log.Printf("Added %d blocks (downloaded %.2f%% of the blockchain). Current bc height: %d\n", len(blocks),
					(float64(blockchain.BcLength)/float64(currBcHeight))*100, blockchain.BcLength)
			}
		}

		// If error
		if !res || !checkpCorresp {
			conn.Close()
			isEstablished, addrInd = conn.EstablishAddresses(node_settings.Settings.OtherNodesAddresses, addrInd, "")
			if !isEstablished {
				return false
			}
		}
	}

	return true
}
