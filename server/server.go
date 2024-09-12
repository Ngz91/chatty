package main

import (
	"chatty/util"
	"context"
	"fmt"
	"log"
	"net"
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

func (s *Server) handleConnection(ctx context.Context, conn net.Conn, buf []byte, broker *util.Broker[string], cArray *[]net.Conn) {
	defer conn.Close()

	msgCh := broker.Subscribe()
	ip := conn.RemoteAddr().String()

	for {
		select {
		case <-ctx.Done():
			broker.Unsubscribe(msgCh)
			return
		case newMsg := <-msgCh:
			log.Println(newMsg)
			for _, c := range *cArray {
				if c != conn {
					log.Printf("Sending from %s to %s: %s", ip, c.RemoteAddr(), newMsg)
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
			msg := string(buf[:n])

			log.Printf("Client %s: %s", ip, msg)
			if msg != "quit" {
				broker.Publish(msg)
			}

			conn.Write([]byte(msg))
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

	defer func() {
		l.Close()
		cancel()
	}()

	b := util.NewBroker[string]() // Used as broadcast channel
	go b.Start()
	defer func() {
		b.Stop()
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
