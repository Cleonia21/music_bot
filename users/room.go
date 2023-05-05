package users

import (
	"errors"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type room struct {
	Bot   *telego.Bot
	Root  *User
	Users map[int64]*User
	Title string
	Pass  string
}

func (r *room) constructor(user *User, title string, pass string) {
	r.Root = user
	r.Users = make(map[int64]*User)
	r.Title = title
	r.Pass = pass
}

func (r *room) out(user *User) {

	delete(r.Users, user.ID)
}

func (r *room) add(user *User, pass string) error {
	if pass != r.Pass {
		return errors.New("incorrect pass")
	}
	r.Users[user.ID] = user
	return nil
}

func (r *room) sendMessage(text string) error {
	params := tu.Message(r.Root.ChatID, text)
	_, err := r.Bot.SendMessage(params)
	if err != nil {
		return err
	}
	return nil
}

func (r *room) sendAudioPackage() error {
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
