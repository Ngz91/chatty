package main

import (
	"bufio"
	"chatty/client/services"
	"chatty/utils"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, 1024)

	cIp := conn.LocalAddr().String()

	ipPort := strings.Split(cIp, ":")
	ip := ipPort[0]
	port := ipPort[1]

	reader := bufio.NewReader(os.Stdin)

	var wg sync.WaitGroup
	wg.Add(1)

	log.Printf("Welcome %s", cIp)
	go services.CheckNewMsg(conn, buf, &wg)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		msg = strings.TrimSpace(msg)

		protoMsg, err := utils.MarshalNewMsg(ip, port, msg)
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
