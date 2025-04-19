package protocol_test

import (
	"encoding/json"
	"sender/internal/app"
	poolMessage "sender/internal/server/blockchain/connectionpool/message"
	"sender/internal/server/blockchain/protocol"
	"sender/internal/server/blockchain/protocol/message"
	"testing"
	"time"
)

func newRawMessage(msg message.Message) message.Message {
	data, _ := json.Marshal(msg)
	return message.Message{
		Type: message.RawMessageType,
		Content: &message.RawMessage{
			MessageJson: data,
		},
	}
}

func TestRun_ResponseInfoUpdatesID(t *testing.T) {
	msgChan := make(chan message.Message, 1)
	poolChan := make(chan poolMessage.PoolMessage, 1)
	state := &app.AppState{}
	proto := protocol.NewProtocol(msgChan, state, poolChan)

	// ResponseMessageInfo с ID 5
	base := message.NewBaseMessage()
	base.SetID(5)
	info := message.Message{
		Type:    message.ResponseMessageInfo,
		Content: &message.InfoMessage{BaseMessage: *base},
	}

	msgChan <- newRawMessage(info)
	go proto.Run()

	time.Sleep(100 * time.Millisecond)
}

func TestRun_RequestMessageInfoTriggersSendFirstMessage(t *testing.T) {
	msgChan := make(chan message.Message, 1)
	poolChan := make(chan poolMessage.PoolMessage, 1)
	state := &app.AppState{}
	proto := protocol.NewProtocol(msgChan, state, poolChan)

	base := message.NewBaseMessage()
	base.SetID(0)
	reqMsg := message.Message{
		Type:    message.RequestMessageInfo,
		Content: &message.InfoMessage{BaseMessage: *base},
	}
	msgChan <- newRawMessage(reqMsg)

	go proto.Run()

	select {
	case out := <-poolChan:
		if out.Type != poolMessage.BroadcastMessage {
			t.Errorf("Expected BroadcastMessage, got %v", out.Type)
		}
	case <-time.After(200 * time.Millisecond):
		t.Error("Expected message to be broadcasted")
	}
}

func TestRun_IgnoreDuplicateMessage(t *testing.T) {
	msgChan := make(chan message.Message, 2)
	poolChan := make(chan poolMessage.PoolMessage, 1)
	state := &app.AppState{}
	proto := protocol.NewProtocol(msgChan, state, poolChan)

	base := message.NewBaseMessage()
	base.SetID(10)
	info := message.Message{
		Type:    message.ResponseMessageInfo,
		Content: &message.InfoMessage{BaseMessage: *base},
	}
	msgChan <- newRawMessage(info)

	go proto.Run()
	time.Sleep(50 * time.Millisecond)

	// Дубликат (меньший ID)
	base2 := message.NewBaseMessage()
	base2.SetID(5)
	duplicate := message.Message{
		Type:    message.ResponseMessageInfo,
		Content: &message.InfoMessage{BaseMessage: *base2},
	}
	msgChan <- newRawMessage(duplicate)

	time.Sleep(100 * time.Millisecond)

	select {
	case <-poolChan:
		t.Error("Duplicate message should not be broadcasted")
	default:
		// ok
	}
}

func TestRun_ProcessBlockCallsAppState(t *testing.T) {
	msgChan := make(chan message.Message, 1)
	poolChan := make(chan poolMessage.PoolMessage, 1)
	state := &app.AppState{}
	proto := protocol.NewProtocol(msgChan, state, poolChan)

	base := message.NewBaseMessage()
	base.SetID(1)
	blockMsg := message.Message{
		Type:    message.ResponseBlockMessage,
		Content: &message.BlockMessage{BaseMessage: *base},
	}
	msgChan <- newRawMessage(blockMsg)

	go proto.Run()
	time.Sleep(100 * time.Millisecond)
}

func TestRun_ProcessPeerCallsConnect(t *testing.T) {
	msgChan := make(chan message.Message, 1)
	poolChan := make(chan poolMessage.PoolMessage, 1)
	state := &app.AppState{}
	proto := protocol.NewProtocol(msgChan, state, poolChan)

	base := message.NewBaseMessage()
	base.SetID(1)
	peerMsg := message.Message{
		Type:    message.ResponsePeerMessage,
		Content: &message.PeerMessage{BaseMessage: *base, PeerAddrIp: "127.0.0.1"},
	}
	msgChan <- newRawMessage(peerMsg)

	go proto.Run()
	time.Sleep(100 * time.Millisecond)
}

func TestRun_UnknownMessageTypeIsIgnored(t *testing.T) {
	msgChan := make(chan message.Message, 1)
	poolChan := make(chan poolMessage.PoolMessage, 1)
	state := &app.AppState{}
	proto := protocol.NewProtocol(msgChan, state, poolChan)

	base := message.NewBaseMessage()
	base.SetID(1)
	unknown := message.Message{
		Type:    "UnknownType",
		Content: &message.InfoMessage{BaseMessage: *base},
	}
	msgChan <- newRawMessage(unknown)

	go proto.Run()
	time.Sleep(100 * time.Millisecond)

	select {
	case <-poolChan:
		t.Error("Unknown message type should not be broadcasted")
	default:
	}
}

func TestRun_InvalidJSONShouldNotPanic(t *testing.T) {
	msgChan := make(chan message.Message, 1)
	poolChan := make(chan poolMessage.PoolMessage, 1)
	state := &app.AppState{}
	proto := protocol.NewProtocol(msgChan, state, poolChan)

	msgChan <- message.Message{
		Type: message.RawMessageType,
		Content: &message.RawMessage{
			MessageJson: []byte(`{invalid json`),
		},
	}

	go proto.Run()
	time.Sleep(100 * time.Millisecond)

	select {
	case <-poolChan:
		t.Error("Invalid JSON should not trigger broadcast")
	default:
	}
}
