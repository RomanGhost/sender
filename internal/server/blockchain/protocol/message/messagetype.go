package message

type MessageType string

const (
	RawMessageType MessageType = "RawMessage"

	RequestMessageInfo         MessageType = "RequestMessageInfo"
	ResponseMessageInfo        MessageType = "ResponseMessageInfo"
	ResponseTransactionMessage MessageType = "ResponseTransactionMessage"
	ResponseBlockMessage       MessageType = "ResponseBlockMessage"
	ResponsePeerMessage        MessageType = "ResponsePeerMessage"
	ResponseTextMessage        MessageType = "ResponseTextMessage"
)
