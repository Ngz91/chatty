package utils

import (
	chat "chatty/proto/v1"

	"google.golang.org/protobuf/proto"
)

func UnmarshalMsg(data []byte) (*chat.ServerMessage, error) {
	newMsg := &chat.ServerMessage{}
	err := proto.Unmarshal(data, newMsg)

	if err != nil {
		return &chat.ServerMessage{}, err
	}
	return newMsg, nil
}
