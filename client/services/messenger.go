package services

import chat "chatty/proto/v1"

type Messenger interface {
	CreateRoomMsg(room string, user *chat.User, roomCreate bool) *chat.RoomMsg
	SendServerMsg(ip string, port string, user *chat.User, msg string) ([]byte, error)
	SendRoomRequest(roomMsg *chat.RoomMsg) error
}
