package message

import "time"

// Интерфейс для всех сообщений
type Message interface {
	GetID() uint64
	SetID(id uint64)
	MessageType() string // Метод для определения типа сообщения
}

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

func (m *BaseMessage) MessageType() string {
	return "BaseMessage"
}
