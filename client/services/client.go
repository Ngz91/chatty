package services

import (
	"chatty/common"
	chat "chatty/proto/v1"
	"net"

	"google.golang.org/protobuf/proto"
)

type Client interface {
	Connect(ip string) error
	Disconnect()
	GetLocalAddr() string
	Write(data []byte)
	Read() ([]byte, error)
}

type TcpClienter interface {
	Client
	Messenger
}

type TcpClient struct {
	conn net.Conn
	buf  []byte
}

func NewTcpClient() TcpClienter {
	return &TcpClient{
		buf: make([]byte, 1500),
	}
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

func (tc *TcpClient) CreateRoomMsg(room string, user *chat.User, roomCreate bool) *chat.RoomMsg {
	roomMsg := chat.RoomMsg{
		Name:   room,
		Create: roomCreate,
		User:   user,
	}
	return &roomMsg
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

func (tc *TcpClient) SendServerMsg(ip string, port string, user *chat.User, msg string) ([]byte, error) {
	m, err := common.MarshalServerMessage(ip, port, user, msg)
	if err != nil {
		return []byte{}, err
	}
	return m, nil
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
