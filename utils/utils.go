package utils

import (
	chat "chatty/proto/v1"

	"google.golang.org/protobuf/proto"
)

func MarshalNewMsg(ip, port, data string) ([]byte, error) {
	msgInfo := chat.Message{
		Ip:      ip,
		Port:    port,
		Content: data,
	}

	protoMsg, err := proto.Marshal(&msgInfo)
	if err != nil {
		return []byte{}, err
	}

	return protoMsg, nil
}

func UnmarshalMsg(data []byte) (*chat.Message, error) {
	newMsg := &chat.Message{}
	err := proto.Unmarshal(data, newMsg)

	if err != nil {
		return &chat.Message{}, err
	}
	return newMsg, nil
}
