package serializemessage

import (
	"encoding/json"
	"sender/internal/jsonutil"
	"sender/internal/server/p2pprotocol/message"
	"sender/internal/server/p2pprotocol/message/request"
	"sender/internal/server/p2pprotocol/message/responce"
)

// Обобщённое сообщение
type GenericMessage struct {
	Type    string          `json:"type"`
	Content message.Message `json:"content"`
}

func NewGenericMessage(m message.Message) *GenericMessage {
	return &GenericMessage{
		Type:    m.MessageType(),
		Content: m,
	}
}

func (gm *GenericMessage) ToJSON() ([]byte, error) {
	return jsonutil.ToJSON(gm)
}

type GenericMessageSerialize struct {
	Type    string          `json:"type"`
	Content json.RawMessage `json:"content"`
}

func FromJSON(jsonByte []byte) (*GenericMessage, error) {
	var gms GenericMessageSerialize
	err := jsonutil.FromJSON(jsonByte, &gms)

	var newMessage message.Message
	switch gms.Type {
	case message.RequestBlocksBeforeMessage.String():
		{
			var bbm request.BlocksBeforeMessage
			newMessage = &bbm
		}
	case message.RequestLastNBlocksMessage.String():
		{
			var lnbm request.BlocksBeforeMessage
			newMessage = &lnbm
		}
	case message.RequestMessageInfo.String():
		{
			var rfm request.InfoMessage
			newMessage = &rfm
		}

	case message.ResponseBlockMessage.String():
		{
			var bm responce.BlockMessage
			newMessage = &bm
		}
	case message.ResponseChainMessage.String():
		{
			var cm responce.ChainMessage
			newMessage = &cm
		}
	case message.ResponseMessageInfo.String():
		{
			var im responce.InfoMessage
			newMessage = &im
		}
	case message.ResponseTransactionMessage.String():
		{
			var tm responce.TransactionMessage
			newMessage = &tm
		}
	case message.ResponsePeerMessage.String():
		{
			var pm responce.PeerMessage
			newMessage = &pm
		}
	}
	jsonutil.FromJSON(gms.Content, &newMessage)

	gm := NewGenericMessage(newMessage)

	return gm, err
}

// Базовая структура для хранения ID и временной метки
