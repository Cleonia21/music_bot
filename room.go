package main

import (
	"container/list"
	"errors"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

//users -> room
//room -> users

type userT struct {
	user   *User
	audios *list.List
}

type Room struct {
	Bot   *telego.Bot
	ID    int64
	Root  *User
	Users map[int64]userT
	Title string
	Pass  string
}

func (r *Room) CreateRoom(rootUser *User, title string, pass string) {
	r.ID = rootUser.ID
	r.Root = rootUser
	r.Users = make(map[int64]userT)
	r.Title = title
	r.Pass = pass
}

func (r *Room) Dell(user *User) {

}

func (r *Room) DellUser(user *User) {
	delete(r.Users, user.ID)
}

func (r *Room) Join(user *User, pass string) error {
	if pass != r.Pass {
		return errors.New("incorrect pass")
	}
	r.Users[user.ID] = userT{user, list.New()}
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

func (r *Room) AddAudio(user *User, params *telego.SendAudioParams) {
	r.Users[user.ID].audios.PushBack(params)
}

func (r *Room) SendAudioPackage() error {
	//for i := r.Users.Front(); i != nil; i = i.Next() {
	//
	//}
	//for _, val := range r.Users {
	//	if val == nil || val.Messages.Len() > 0 {
	//		continue
	//	}
	//	_, err := r.Bot.SendAudio(val.Messages.Front().Value.(*telego.SendAudioParams))
	//	if err != nil {
	//		return err
	//	}
	//	val.Messages.Remove(val.Messages.Front())
	//}
	return nil
}
