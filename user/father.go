package user

import (
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
)

type userFather struct {
	tg *telego.Bot
	id telego.ChatID
}

func (u *userFather) sendText(text string) {
	if text == "" {
		return
	}
	_, err := u.tg.SendMessage(telegoutil.Message(u.id, text))
	if err != nil {
		u.tg.Logger().Errorf(err.Error())
	}
}

func (u *userFather) sendMessage(msg *telego.SendMessageParams) {
	_, err := u.tg.SendMessage(msg)
	if err != nil {
		u.tg.Logger().Errorf(err.Error())
	}
}

func (u *userFather) getID() telego.ChatID {
	return u.id
}
