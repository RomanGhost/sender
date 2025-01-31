package blockchain

import (
	"bufio"
	"bytes"
	"net"
	"time"
)

type MockConn struct {
	Reader *bufio.Reader
	Writer *bytes.Buffer
}

func NewMockConn() *MockConn {
	return &MockConn{
		Reader: bufio.NewReader(bytes.NewReader([]byte{})),
		Writer: new(bytes.Buffer),
	}
}

func (mc *MockConn) Read(b []byte) (n int, err error) {
	return mc.Reader.Read(b)
}

func (mc *MockConn) Write(b []byte) (n int, err error) {
	return mc.Writer.Write(b)
}

func (mc *MockConn) Close() error {
	return nil
}

func (mc *MockConn) LocalAddr() net.Addr {
	return &net.TCPAddr{}
}

func (mc *MockConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{}
}

func (mc *MockConn) SetDeadline(t time.Time) error {
	return nil
}

func (mc *MockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (mc *MockConn) SetWriteDeadline(t time.Time) error {
	return nil
}
