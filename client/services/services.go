package services

import (
	chat "chatty/proto/v1"
	"context"
	"io"
	"log"
	"sync"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

func NewUser(username string) *chat.User {
	u := &chat.User{
		Id:       uuid.NewString(),
		Username: username,
	}
	return u
}

func CheckNewMsg(ctx context.Context, c Client, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			data, err := c.Read()
			if err != nil {
				if err != io.EOF {
					log.Println("Connection closed")
					return
				}
				return
			}
			m := &chat.ClientMessage{}
			err = proto.Unmarshal(data, m)

			u := m.GetUser().GetUsername()
			if u == "" {
				u = m.GetUser().GetId()
			}
			log.Printf("Client %s says: %s", u, m.GetContent())
		}
	}
}

func CheckStatus(c Client) bool {
	data, err := c.Read()
	if err != nil {
		if err == io.EOF {
			log.Fatal("Connection closed")
		}
		log.Fatal(err)
	}
	o := &chat.Operation{}
	err = proto.Unmarshal(data, o)

	if err != nil {
		log.Fatal(err)
	}

	if o.Success == 1 {
		return true
	}
	return false
}
