package services

import (
	chat "chatty/proto/v1"
	"net"

	"google.golang.org/protobuf/proto"
)

type Client interface {
	Connect(ip string) error
	Disconnect()
	SendRoomRequest(roomMsg *chat.RoomMsg) error
	GetLocalAddr() string
	Write(data []byte)
	Read() ([]byte, error)
}

type TcpClient struct {
	conn net.Conn
	buf  []byte
}

func (tc *TcpClient) Connect(ip string) error {
	c, err := net.Dial("tcp", ip)
	if err != nil {
		return err
	}
	tc.conn = c
	return nil
}

func (tc *TcpClient) Disconnect() {
	tc.conn.Close()
}

func (tc *TcpClient) SendRoomRequest(roomMsg *chat.RoomMsg) error {
	// Send a room request to the server
	// The server handles the create/join room logic
	rMsg, err := proto.Marshal(roomMsg)
	if err != nil {
		return err
	}
	tc.Write(rMsg)
	return nil
}

func (tc *TcpClient) GetLocalAddr() string {
	return tc.conn.LocalAddr().String()
}

func (tc *TcpClient) Write(data []byte) {
	tc.conn.Write(data)
}

func (tc *TcpClient) Read() ([]byte, error) {
	n, err := tc.conn.Read(tc.buf)
	if err != nil {
		return []byte{}, err
	}
	data := tc.buf[:n]
	return data, nil
}

func NewTcpClient() Client {
	return &TcpClient{
		buf: make([]byte, 1500),
	}
}
