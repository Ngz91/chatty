package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

func checkNewMsg(conn net.Conn, buf []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	// TODO Check if the message is not from the same user (Probably create a is the server a struct to hold the info necessary like a uuid, channel etc)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println("Connection closed")
				return
			}
			return
		}
		log.Printf("Received from server: %s", string(buf[:n]))
	}
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, 1024)
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

		conn.Write([]byte(msg))

		if msg == "quit" {
			conn.Close()
			break
		}
	}
	wg.Wait()
}
