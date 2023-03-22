package synchronization

import (
	"bavovnacoin/blockchain"
	"bavovnacoin/networking"
	"bavovnacoin/node_controller/node_settings"
	"log"
)

func StartInitSync(printLog bool) bool {
	var conn networking.Connection
	var addrInd int = 0

	var isEstablished bool
	isEstablished, addrInd = conn.EstablishAddresses(node_settings.Settings.OtherNodesAddresses, addrInd)
	if !isEstablished {
		return false
	}

	InitCheckpoints()
	var blockReqInd uint64 = 0
	checkpCorresp := false
	for true {
		blocks, res := conn.RequestBlocks(blockReqInd)
		if res {
			if len(blocks) == 0 {
				return true
			}

			for i := 0; i < len(blocks); i++ {
				checkpCorresp = checkForCheckpCorrespondence(blockReqInd, blocks[i])
				if checkpCorresp {
					blockchain.AddBlockToBlockchain(blocks[i])
					blockReqInd++
				} else {
					if printLog {
						log.Printf("Address %s sent an incorrect block. Selecting next address\n",
							node_settings.Settings.OtherNodesAddresses[addrInd])
					}
					break
				}
			}
		}

		// If error
		if !res || !checkpCorresp {
			conn.Close()
			isEstablished, addrInd = conn.EstablishAddresses(node_settings.Settings.OtherNodesAddresses, addrInd)
			if !isEstablished {
				return false
			}
		}
	}
	return true
}
