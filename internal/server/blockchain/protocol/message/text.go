package message

type TextMessage struct {
	BaseMessage
	Message string `json:"message"`
}
