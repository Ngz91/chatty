package main

import (
	"bufio"
	"chatty/client/services"
	"flag"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

var username string

func main() {
	flag.StringVar(&username, "user", "", "Username to be used.")
	// room := flag.String("room", "", "Room to enter. If create_room is set to true the name of the room will be used to create a new room.")
	// createRoom := flag.Bool("create_room", false, "Create a new room based on the provided room name.")

	flag.Parse()

	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, 1024)

	cIp := conn.LocalAddr().String()

	user := services.NewUser(conn, username)

	if username == "" {
		log.Printf("No username provided, using uuid %s", user.GetId())
	}

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

		protoMsg, err := services.NewServerMessage(ip, port, user, msg)
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
