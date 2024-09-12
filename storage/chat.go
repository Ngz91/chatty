package storage

import (
	chat "chatty/proto/v1"
	"io"

	"google.golang.org/protobuf/proto"
)

func Load(r io.Reader) ([]*chat.Message, error) {
	// TODO Loading based on room
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var chores chat.Messages
	return chores.Messages, proto.Unmarshal(b, &chores)
}

func Flush(w io.Writer, messages []*chat.Message) error {
	// Flush the the messages to a File
	b, err := proto.Marshal(&chat.Messages{Messages: messages})
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}
