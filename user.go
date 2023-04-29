package main

import (
	"container/list"
	"github.com/mymmrac/telego"
)

type User struct {
	ChatID     telego.ChatID
	ID         int64
	Room       *Room
	isRoomRoot bool
	Nik        string
	Messages   *list.List //[]*telego.SendAudioParams

}

func CreateUser(ChatID telego.ChatID, ID int64, nik string) *User {
	u := User{
		ChatID: ChatID,
		ID:     ID,
		Nik:    nik,
	}
	u.Messages = list.New()
	return &u
}

func (u *User) addMessage(params *telego.SendAudioParams) {
	u.Messages.PushBack(params)
}

func (u *User) connectToRoom(room *Room, isRoot bool) {
	u.Room = room
	u.isRoomRoot = isRoot
}

