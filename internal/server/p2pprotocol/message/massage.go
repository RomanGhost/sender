package message

import (
	"sender/internal/jsonutil"
	"time"
)

// Интерфейс для всех сообщений
type Message interface {
	GetID() uint64
	SetID(id uint64)
	MessageType() string // Метод для определения типа сообщения
}

// Обобщённое сообщение
type GenericMessage struct {
	Type    string  `json:"type"`
	Content Message `json:"content"`
}

func NewGenericMessage(m Message) *GenericMessage {
	return &GenericMessage{
		Type:    m.MessageType(),
		Content: m,
	}
}

func (gm *GenericMessage) ToJSON() ([]byte, error) {
	return jsonutil.ToJSON(gm)
}

func FromJSON(jsonByte []byte) (*GenericMessage, error) {
	var gm GenericMessage
	err := jsonutil.FromJSON(jsonByte, gm)
	return &gm, err
}

// Базовая структура для хранения ID и временной метки
type BaseMessage struct {
	ID        uint64    `json:"id"`
	Timestamp time.Time `json:"time_stamp"`
}

func (m *BaseMessage) GetID() uint64 {
	return m.ID
}

func (m *BaseMessage) SetID(id uint64) {
	m.ID = id
}
