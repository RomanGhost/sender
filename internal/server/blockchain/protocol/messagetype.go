package protocol

type MessageType string

const (
	RawMessage MessageType = "RawMessage"

	RequestMessageInfo         MessageType = "RequestMessageInfo"
	ResponseMessageInfo        MessageType = "ResponseMessageInfo"
	ResponseTransactionMessage MessageType = "ResponseTransactionMessage"
	ResponseBlockMessage       MessageType = "ResponseBlockMessage"
	ResponsePeerMessage        MessageType = "ResponsePeerMessage"
	ResponseTextMessage        MessageType = "ResponseTextMessage"
)
