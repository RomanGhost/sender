package message

type MessageType string

const (
	RequestBlocksBeforeMessage MessageType = "RequestBlocksBeforeMessage"
	RequestMessageInfo         MessageType = "RequestMessageInfo"
	RequestLastNBlocksMessage  MessageType = "RequestLastNBlocksMessage"

	ResponseBlockMessage       MessageType = "ResponseBlockMessage"
	ResponseChainMessage       MessageType = "ResponseChainMessage"
	ResponseMessageInfo        MessageType = "ResponseMessageInfo"
	ResponseTransactionMessage MessageType = "ResponseTransactionMessage"
)

func (mt MessageType) String() string {
	return string(mt)
}
