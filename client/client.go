package main

import (
	"bufio"
	"chatty/client/services"
	"context"
	"flag"
	"io"
	"log"
	"os"
	"strings"
	"sync"
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

	ctx, cancel := context.WithCancel(context.Background())

	client := services.NewTcpClient()
	err := client.Connect("127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}

	cIp := client.GetLocalAddr()
	user := services.NewUser(username)

	if user.GetUsername() == "" {
		log.Printf("No username provided, using uuid %s", user.GetId())
	}

	ipPort := strings.Split(cIp, ":")
	ip := ipPort[0]
	port := ipPort[1]

	roomMsg := client.CreateRoomMsg(room, user, roomCreate)

	err = client.SendRoomRequest(roomMsg)
	if err != nil {
		client.Disconnect()
		log.Fatal(err)
	}

	// Confirm that the room was created or exists and the user was added
	s := services.CheckStatus(client)
	if s == false {
		log.Fatal("Server could not create/join room")
		client.Disconnect()
	}

	reader := bufio.NewReader(os.Stdin)

	var wg sync.WaitGroup
	wg.Add(1)

	go services.CheckNewMsg(ctx, client, &wg)

	defer func() {
		cancel()
		client.Disconnect()
	}()

	log.Printf("Welcome %s to room %s", cIp, room)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			log.Fatal("Forcefully exited the program")
		}
		msg = strings.TrimSpace(msg)

		protoMsg, err := client.SendServerMsg(ip, port, user, msg)
		if err != nil {
			log.Fatal(err)
		}

		client.Write(protoMsg)

		if msg == "/quit" {
			cancel()
			client.Disconnect()
			break
		}
	}
	wg.Wait()
}
