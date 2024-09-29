package main

import (
	"bufio"
	"chatty/client/services"
	chat "chatty/proto/v1"
	"flag"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"google.golang.org/protobuf/proto"
)

var username string
var joinRoom string
var createRoom string

func main() {
	flag.StringVar(&username, "user", "", "Username to be used.")
	flag.StringVar(&joinRoom, "join_room", "", "Room to join. If create_room is set to true the name of the room will be used to create a new room.")
	flag.StringVar(&createRoom, "create_room", "", "Create a new room based on the provided room name.")

	flag.Parse()

	var room string
	var roomCreate bool

	if joinRoom == "" && createRoom == "" {
		log.Fatal("Provide name of the room to join or create...")
	}

	if joinRoom != "" && createRoom != "" {
		log.Fatal("To create and join a room use only create_room...")
	}

	if joinRoom == "" && createRoom != "" {
		room = createRoom
		roomCreate = true
	} else {
		room = joinRoom
		roomCreate = false
	}

	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, 1500)

	cIp := conn.LocalAddr().String()

	user := services.NewUser(conn, username)

	if user.GetUsername() == "" {
		log.Printf("No username provided, using uuid %s", user.GetId())
	}

	ipPort := strings.Split(cIp, ":")
	ip := ipPort[0]
	port := ipPort[1]

	roomMsg := chat.RoomMsg{
		Name:   room,
		Create: roomCreate,
		User:   user,
	}

	// Send a room request to the server
	// The server handles the create/join room logic
	rMsg, err := proto.Marshal(&roomMsg)
	if err != nil {
		log.Fatal(err)
	}
	conn.Write(rMsg)

	// Confirm that the room was created or exists and the user was added
	s := services.CheckStatus(conn, buf)
	if s == false {
		log.Fatal("Server could not create/join room")
		conn.Close()
	}

	reader := bufio.NewReader(os.Stdin)

	var wg sync.WaitGroup
	wg.Add(1)

	log.Printf("Welcome %s to room %s", cIp, room)
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
