package message

import "time"

type BaseMessage struct {
	ID        uint64 `json:"id"`
	TimeStamp int64  `json:"time_stamp"`
}

func NewBaseMessage() *BaseMessage {
	return &BaseMessage{
		ID:        0,
		TimeStamp: time.Now().UTC().Unix(),
	}
}

func (bm *BaseMessage) GetID() uint64 {
	return bm.ID
}

func (bm *BaseMessage) SetID(newID uint64) {
	bm.ID = newID
}

func (bm *BaseMessage) GetTime() int64 {
	return bm.TimeStamp
}
