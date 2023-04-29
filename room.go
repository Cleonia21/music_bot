package main

import (
	//"container/list"
	"errors"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

//users -> room
//room -> users

type Room struct {
	Bot   *telego.Bot
	ID    int64
	Root  *User
	Users []*User //user_id / nik
	Title string
	Pass  string
}

func CreateRoom(rootUser *User, title string, pass string) *Room {
	var r Room
	r.ID = rootUser.ID
	r.Root = rootUser
	r.Users = make([]*User, 20)
	r.Title = title
	r.Pass = pass
	return &r
}

func (r *Room) DellUser(user *User) {
	for i, val := range r.Users {
		if val == user {
			r.Users[i] = nil
		}
	}
}

func (r *Room) Join(user *User, pass string) error {
	if pass != r.Pass {
		return errors.New("incorrect pass")
	}
	r.Users[user.ID] = user
	return nil
}

func (r *Room) SendMessage(text string) error {
	params := tu.Message(r.Root.ChatID, text)
	_, err := r.Bot.SendMessage(params)
	if err != nil {
		return err
	}
	return nil
}

func (r *Room) SendAudioPackage() error {
	for _, val := range r.Users {
		if val == nil || val.Messages.Len() > 0 {
			continue
		}
		_, err := r.Bot.SendAudio(val.Messages.Front().Value.(*telego.SendAudioParams))
		if err != nil {
			return err
		}
		val.Messages.Remove(val.Messages.Front())
	}
	return nil
}
