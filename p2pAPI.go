package rpc

func GetPeerCount() (int, error) {
	peerCount := int(0)
	err := ClientCall("p2p_peerCount", &peerCount)

	return peerCount, err
}

func AddPeer(url string) (uint64, error) {
	peerCount := uint64(0)
	err := ClientCall("p2p_addPeer", &peerCount, url)

	return peerCount, err
}

func RemovePeer(url string) (bool, error) {
	bSuccess := true
	err := ClientCall("p2p_removePeer", &bSuccess, url)

	return bSuccess, err
}

func AddTrustedPeer(url string) (uint64, error) {
	peerCount := uint64(0)
	err := ClientCall("p2p_addTrustedPeer", &peerCount, url)

	return peerCount, err
}

func RemoveTrustedPeer(url string) (bool, error) {
	bSuccess := true
	err := ClientCall("p2p_removeTrustedPeer", &bSuccess, url)

	return bSuccess, err
}

func SelfNode(url string) ([]string, error) {
	selfNodes := []string{}
	err := ClientCall("p2p_selfNode", &selfNodes, url)

	return selfNodes, err
}

func Peers(url string) ([]string, error) {
	peers := []string{}
	err := ClientCall("p2p_peers", &peers, url)

	return peers, err
}



