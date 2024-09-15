package services

import (
	chat "chatty/proto/v1"
	"io"
	"log"
	"net"
	"sync"

	"google.golang.org/protobuf/proto"
)

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
		m := &chat.Message{}
		err = proto.Unmarshal(data, m)
		log.Printf("Client %s in port %s says: %s", m.GetIp(), m.GetPort(), m.GetContent())
	}
}
