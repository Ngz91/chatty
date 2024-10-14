package services

import (
	chat "chatty/proto/v1"
	"io"
	"log"
	"net"
	"sync"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

func NewUser(conn net.Conn, username string) *chat.User {
	u := &chat.User{
		Id:       uuid.NewString(),
		Username: username,
	}
	return u
}

func NewServerMessage(ip string, port string, user *chat.User, msg string) ([]byte, error) {
	sMsg := &chat.ServerMessage{
		Ip:      ip,
		Port:    port,
		User:    user,
		Content: msg,
	}

	protoMsg, err := proto.Marshal(sMsg)
	if err != nil {
		return []byte{}, err
	}

	return protoMsg, nil
}

func CheckNewMsg(conn net.Conn, buf []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println("Connection closed")
				return
			}
			return
		}
		data := buf[:n]
		m := &chat.ClientMessage{}
		err = proto.Unmarshal(data, m)

		u := m.GetUser().GetUsername()
		if u == "" {
			u = m.GetUser().GetId()
		}
		log.Printf("Client %s says: %s", u, m.GetContent())
	}
}

func CheckStatus(conn net.Conn, buf []byte) bool {
	n, err := conn.Read(buf)
	if err != nil {
		if err == io.EOF {
			log.Fatal("Connection closed")
		}
		log.Fatal(err)
	}
	data := buf[:n]
	o := &chat.Operation{}
	err = proto.Unmarshal(data, o)

	if err != nil {
		log.Fatal(err)
	}

	if o.Success == 1 {
		return true
	}
	return false
}
