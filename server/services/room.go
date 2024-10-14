package services

import (
	"net"
)

type Room struct {
	Id    string
	Users []net.Conn
}

func newRoom(id string) *Room {
	return &Room{
		Id:    id,
		Users: make([]net.Conn, 0),
	}
}

func getUsers(r Room) []net.Conn {
	return r.Users
}

func removeUser(id string, conn net.Conn, r Room) Room {
	room := Room{Id: id}
	result := []net.Conn{}
	for _, c := range r.Users {
		if c != conn {
			result = append(result, c)
		}
	}
	room.Users = result
	return room
}

func deleteRoom(id string, rMap map[string]Room) {
	if len(getUsers(rMap[id])) == 0 {
		delete(rMap, id)
	}
}
