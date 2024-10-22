package services

import (
	"chatty/common"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	host string
	port string
}

type Config struct {
	Host string
	Port string
}

func NewServer(config *Config) *Server {
	return &Server{
		host: config.Host,
		port: config.Port,
	}
}

func (s *Server) Run() {
	ctx, cancel := context.WithCancel(context.Background())

	buf := make([]byte, 1500)
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.host, s.port))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening on address %s", l.Addr())

	b := make(chan []byte, 2)     // Used as broadcast channel
	rMap := make(map[string]Room) // Map of rooms created (Room name -> Room struct)

	defer func() {
		l.Close()
		cancel()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		// TODO Handle the deletion of the rooms and close of connections
		// notify the clients of the server shutting down
		l.Close()
		cancel()
		os.Exit(1)
	}()

	for {
		c, err := l.Accept()
		log.Printf("New connection %s", c.RemoteAddr().String())
		if err != nil {
			log.Fatal(err)
			c.Close()
		}

		rName, err := s.createJoinRoom(c, buf, rMap)
		if err != nil {
			log.Printf("Error creating/joining room, closing connection from %s", c.RemoteAddr())
			log.Println(err)
			c.Close()
			continue
		}
		go s.handleConnection(ctx, c, buf, b, rMap, rName)
	}
}

func (s *Server) createJoinRoom(conn net.Conn, buf []byte, rMap map[string]Room) (string, error) {
	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("Client %s closed connection", conn.RemoteAddr().String())
		conn.Close()
	}
	data := buf[:n]
	roomMsg, err := common.UnmarshalRoomMsg(data)
	if err != nil {
		log.Fatal(err)
	}

	var oStatus []byte // Used to inform the client of the creation/join of a room

	roomName := roomMsg.GetName()

	if roomMsg.GetCreate() == true {
		log.Printf("Creating new room: %s", roomMsg.GetName())
		room := *newRoom(roomName)
		room.Users = append(room.Users, conn) // Add user to room
		rMap[roomName] = room
		oStatus, err = common.MarshalStatusMsg(1)
		if err != nil {
			log.Fatal("unexpected error encoding status message when creating a room")
		}
		conn.Write(oStatus)
	} else {
		if r, ok := rMap[roomName]; !ok {
			oStatus, err = common.MarshalStatusMsg(2)
			if err != nil {
				log.Fatal("unexpected error encoding status message when creating a room with the same name")
			}
			conn.Write(oStatus)
			errString := fmt.Sprintf("No room named %s exists.", roomMsg.Name)
			return "", errors.New(errString)
		} else {
			log.Printf("Joining %s to room %s", roomMsg.GetUser().GetId(), roomName)
			r.Users = append(r.Users, conn)
			rMap[roomName] = r // Update Room
			oStatus, err = common.MarshalStatusMsg(1)
			if err != nil {
				log.Fatal("unexpected error encoding status message when joining an existing room")
			}
			conn.Write(oStatus)
		}
	}
	return roomName, nil
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn, buf []byte, bCh chan []byte, rMap map[string]Room, rName string) {
	defer conn.Close()

	ip := conn.RemoteAddr().String()

	for {
		select {
		case <-ctx.Done():
			return
		case newMsg := <-bCh:
			for _, c := range getUsers(rMap[rName]) {
				if c != conn {
					c.Write([]byte(newMsg))
				}
			}
		default:
			n, err := conn.Read(buf)
			if err != nil {
				log.Printf("Client %s closed connection", conn.RemoteAddr().String())
				conn.Close()
				r := removeUser(rName, conn, rMap[rName])
				rMap[rName] = r

				deleteRoom(rName, rMap)

				return
			}
			data := buf[:n]

			newMsg, err := common.UnmarshalMsg(data)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("Client %s: %s", ip, newMsg.GetContent())
			if newMsg.Content != "quit" {
				m, err := common.MarshalClientMsg(newMsg)
				if err != nil {
					log.Fatal(err)
				}
				bCh <- m
			}
		}
	}
}
