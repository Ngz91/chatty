package common

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

func UnmarshalRoomMsg(data []byte) (*chat.RoomMsg, error) {
	roomMsg := &chat.RoomMsg{}
	err := proto.Unmarshal(data, roomMsg)

	if err != nil {
		return &chat.RoomMsg{}, err
	}
	return roomMsg, nil
}

func MarshalStatusMsg(status uint32) ([]byte, error) {
	opMsg := &chat.Operation{
		Success: status,
	}
	o, err := proto.Marshal(opMsg)
	if err != nil {
		return []byte{}, err
	}
	return o, nil
}

func MarshalServerMessage(ip string, port string, user *chat.User, msg string) ([]byte, error) {
	sMsg := &chat.ServerMessage{
		Ip:      ip,
		Port:    port,
		User:    user,
		Content: msg,
	}

	protoMsg, err := proto.Marshal(sMsg)
	if err != nil {
		return []byte{}, err
	}

	return protoMsg, nil
}
