package message

type PeerMessage struct {
	BaseMessage
	PeerAddrIp string `json:"peer_address"`
}
