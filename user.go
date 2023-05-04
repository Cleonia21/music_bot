package main

import (
	"github.com/mymmrac/telego"
)

const ()

type User struct {
	ChatID      telego.ChatID
	ID          int64
	Room        *Room
	isRoomRoot  bool
	Nik         string
	MessHistory map[int]string
	State       int
}

func CreateUser(ChatID telego.ChatID, ID int64, nik string) *User {
	u := User{
		ChatID:      ChatID,
		ID:          ID,
		Nik:         nik,
		MessHistory: make(map[int]string),
	}
	return &u
}

func (u *User) addAudio(params *telego.SendAudioParams) {
	u.Room.AddAudio(u, params)
}

func (u *User) joinRoom(room *Room, isRoot bool) {
	u.Room = room
	u.isRoomRoot = isRoot
}
