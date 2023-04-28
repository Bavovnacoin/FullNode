package networking_p2p

func TryAddCameBlock() {

}

// TODO: launch with 'go'
func SendOutNewBlock() {
	peerIds := Peer.Peerstore().Peers()
	for i := 0; i < len(peerIds); i++ {
		println(peerIds[i].Pretty())
	}
}
