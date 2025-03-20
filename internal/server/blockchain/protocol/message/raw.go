package message

import "encoding/json"

type RawMessage struct {
	BaseMessage
	MessageJson json.RawMessage
}
