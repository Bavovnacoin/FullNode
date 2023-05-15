package networking_p2p

import (
	"bavovnacoin/byteArr"
	"bavovnacoin/node/node_settings"
	"fmt"
	"strings"

	"github.com/libp2p/go-libp2p/core/peer"
)

const ADDR_SEND_AMMOUNT = 10

var IsAddressesRequested bool
var AddrDownloadCounter int

type NodeAddressesResponse struct {
	IsMore    bool
	Ind       int
	Addresses [][]string
}

func (pd *PeerData) TryHandleAddressRequest(data []byte, peerId peer.ID) bool {
	if data[0] == 9 {
		var startFrom int
		byteArr.FromByteArr(data[1:], &startFrom)

		var response NodeAddressesResponse

		for ; startFrom < len(node_settings.Settings.OtherNodesAddresses); startFrom++ {
			if len(response.Addresses) == ADDR_SEND_AMMOUNT {
				response.IsMore = true
				break
			}

			response.Addresses = append(response.Addresses, node_settings.Settings.OtherNodesAddresses[startFrom])
		}

		response.Ind = startFrom

		respByte, _ := byteArr.ToByteArr(response)
		return pd.SendDataOnPeerId(append([]byte{10}, respByte...), peerId)
	} else if data[0] == 10 {
		if IsAddressesRequested {
			var response NodeAddressesResponse
			byteArr.FromByteArr(data[1:], &response)

			for i := 0; i < len(response.Addresses); i++ {
				if fmt.Sprintf("%s/%s", Peer.Peer.Addrs()[0], Peer.Peer.ID().Pretty()) != response.Addresses[i][0] &&
					node_settings.Settings.AddAddress(response.Addresses[i]) {
					AddrDownloadCounter++
					println("Added address: ", strings.Join(response.Addresses[i], "_"))
				}
			}

			node_settings.Settings.WriteSettings()

			if response.IsMore {
				pd.RequestForNodeAddresses(response.Ind + len(response.Addresses))
			} else {
				IsAddressesRequested = false
			}

			return true
		}
	}
	return false
}

func (pd *PeerData) RequestForNodeAddresses(startFrom int) bool {
	startBytes, _ := byteArr.ToByteArr(startFrom)
	return pd.SendDataToAllConnectedPeers(append([]byte{9}, startBytes...))
}
