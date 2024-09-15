package main

import (
	"chatty/utils"
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

func NewServer(config *Config) *Server {
	return &Server{
		host: config.Host,
		port: config.Port,
	}
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn, buf []byte, bCh chan []byte, cArray *[]net.Conn) {
	defer conn.Close()

	ip := conn.RemoteAddr().String()

	for {
		select {
		case <-ctx.Done():
			return
		case newMsg := <-bCh:
			for _, c := range *cArray {
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

			newMsg, err := utils.UnmarshalMsg(data)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("Client %s: %s", ip, newMsg.Content)
			if newMsg.Content != "quit" {
				m, err := proto.Marshal(newMsg)
				if err != nil {
					log.Fatal(err)
				}
				bCh <- m
			}
		}
	}
}

func (s *Server) run() {
	ctx, cancel := context.WithCancel(context.Background())

	buf := make([]byte, 1024)
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

	// TODO Create map of connections and map them to rooms
	var cArray []net.Conn

	for {
		c, err := l.Accept()
		log.Printf("New connection %s", c.RemoteAddr().String())
		if err != nil {
			log.Fatal(err)
			c.Close()
		}
		// TODO Create a room
		cArray = append(cArray, c)
		go s.handleConnection(ctx, c, buf, b, &cArray)
	}
}

func main() {
	s := NewServer(&Config{
		Host: "localhost",
		Port: "8080",
	})
	s.run()
}
