package main

import (
	"bufio"
	chat "chatty/proto/v1"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"google.golang.org/protobuf/proto"
)

func checkNewMsg(conn net.Conn, buf []byte, wg *sync.WaitGroup) {
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
		m := &chat.Message{}
		err = proto.Unmarshal(buf[:n], m)
		log.Printf("Client %s in port %s says: %s", m.GetIp(), m.GetPort(), m.GetContent())
	}
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, 1024)
	ipPort := strings.Split(conn.LocalAddr().String(), ":")
	ip := ipPort[0]
	port := ipPort[1]

	reader := bufio.NewReader(os.Stdin)

	var wg sync.WaitGroup
	wg.Add(1)

	log.Println("-->")
	go checkNewMsg(conn, buf, &wg)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		msg = strings.TrimSpace(msg)

		msgInfo := chat.Message{
			Ip:      ip,
			Port:    port,
			Content: msg,
		}

		protoMsg, err := proto.Marshal(&msgInfo)
		if err != nil {
			log.Fatal(err)
		}

		conn.Write(protoMsg)

		if msg == "quit" {
			conn.Close()
			break
		}
	}
	wg.Wait()
}
