package services

import chat "chatty/proto/v1"

type Messenger interface {
	CreateRoomMsg(room string, user *chat.User, roomCreate bool) *chat.RoomMsg
	SendRoomRequest(roomMsg *chat.RoomMsg) error
}
