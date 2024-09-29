package services

import (
	chat "chatty/proto/v1"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/protobuf/proto"
)

type Server struct {
	host string
	port string
}

type Config struct {
	Host string
	Port string
}

type Room struct {
	Id    string
	Users []net.Conn // Array of connections in the room (use of a map would be more optimal)
}

func NewServer(config *Config) *Server {
	return &Server{
		host: config.Host,
		port: config.Port,
	}
}

func newRoom(id string) *Room {
	return &Room{
		Id: id,
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

	b := make(chan []byte, 2) // Used as broadcast channel

	defer func() {
		l.Close()
		cancel()
	}()

	rMap := make(map[string]Room) // Map of rooms created (Room name -> Room struct)

	for {
		c, err := l.Accept()
		log.Printf("New connection %s", c.RemoteAddr().String())
		if err != nil {
			log.Fatal(err)
			c.Close()
		}

		go s.handleConnection(ctx, c, buf, b, rMap)
	}
}

func (s *Server) createJoinRoom(conn net.Conn, buf []byte, rMap map[string]Room) string {
	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("Client %s closed connection", conn.RemoteAddr().String())
		conn.Close()
	}
	data := buf[:n]
	roomMsg, err := unmarshalRoomMsg(data)
	if err != nil {
		log.Fatal(err)
	}

	var oStatus []byte // Used to inform the client of the creation/join of a room

	r := roomMsg.Name

	if roomMsg.GetCreate() == true {
		log.Printf("Creating new room: %s", roomMsg.GetName())
		room := *newRoom(roomMsg.GetName())
		room.Users = append(room.Users, conn) // add user to room
		rMap[roomMsg.GetName()] = room
		oStatus, err = marshalStatusMsg(true)
		if err != nil {
			log.Fatal("unexpected error encoding status message when creating a room")
		}
		conn.Write(oStatus)
	} else {
		if r, ok := rMap[roomMsg.GetName()]; !ok {
			oStatus, err = marshalStatusMsg(false)
			if err != nil {
				log.Fatal("unexpected error encoding status message when creating a room with the same name")
			}
			conn.Write(oStatus)
		} else {
			log.Printf("Joining %s to room %s", roomMsg.GetUser().GetId(), roomMsg.GetName())
			r.Users = append(r.Users, conn)
			log.Println(r.Users)
			oStatus, err = marshalStatusMsg(true)
			if err != nil {
				log.Fatal("unexpected error encoding status message when joining an existing room")
			}
			conn.Write(oStatus)
		}
	}
	return r
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn, buf []byte, bCh chan []byte, rMap map[string]Room) {
	defer conn.Close()

	ip := conn.RemoteAddr().String()

	rName := s.createJoinRoom(conn, buf, rMap)

	for {
		select {
		case <-ctx.Done():
			return
		case newMsg := <-bCh:
			log.Println(rMap)
			for _, c := range rMap[rName].Users {
				if c != conn {
					c.Write([]byte(newMsg))
				}
			}
		default:
			n, err := conn.Read(buf)
			if err != nil {
				log.Printf("Client %s closed connection", conn.RemoteAddr().String())
				conn.Close()
				return
			}
			data := buf[:n]

			newMsg, err := unmarshalMsg(data)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("Client %s: %s", ip, newMsg.GetContent())
			if newMsg.Content != "quit" {
				m := &chat.ClientMessage{
					User:    newMsg.GetUser(),
					Content: newMsg.GetContent(),
				}
				newMsgProto, err := proto.Marshal(m)
				if err != nil {
					log.Fatal(err)
				}
				bCh <- newMsgProto
			}
		}
	}
}

func unmarshalMsg(data []byte) (*chat.ServerMessage, error) {
	newMsg := &chat.ServerMessage{}
	err := proto.Unmarshal(data, newMsg)

	if err != nil {
		return &chat.ServerMessage{}, err
	}
	return newMsg, nil
}

func unmarshalRoomMsg(data []byte) (*chat.RoomMsg, error) {
	roomMsg := &chat.RoomMsg{}
	err := proto.Unmarshal(data, roomMsg)

	if err != nil {
		return &chat.RoomMsg{}, err
	}
	return roomMsg, nil
}

func marshalStatusMsg(status bool) ([]byte, error) {
	opMsg := &chat.Operation{
		Success: status,
	}
	o, err := proto.Marshal(opMsg)
	if err != nil {
		return []byte{}, err
	}
	return o, nil
}
