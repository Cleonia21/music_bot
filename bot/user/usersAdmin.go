package user

import (
	"MusicBot/audio"
	"MusicBot/log"
	"MusicBot/user/utils"
	utils2 "MusicBot/utils"
	"github.com/mymmrac/telego"
)

type users interface {
	handler(update *telego.Update) (users, bool)
	getID() utils.UserID
}

type Bot interface {
	SendMessage(params *telego.SendMessageParams) (*telego.Message, error)
	SendAudio(params *telego.SendAudioParams) (*telego.Message, error)
	AnswerCallbackQuery(params *telego.AnswerCallbackQueryParams) error
	EditMessageText(params *telego.EditMessageTextParams) (*telego.Message, error)
}

type Admin struct {
	users map[utils.UserID]users
	audio *audio.Audio
}

func Init() *Admin {
	a := Admin{
		users: make(map[utils.UserID]users),
		audio: audio.Init(),
	}
	return &a
}

func (a *Admin) Handler(update *telego.Update) {
	log.Logger.Debugf("get update: " + utils2.UpdateToStr(update))

	id := utils.UpdateToID(update)
	user, ok := a.users[id]

	if !ok {
		user := &unregUser{}
		user.Init(
			id,
		)
		a.users[id] = users(user)
	} else {
		newUser, needInit := user.handler(update)
		if needInit {
			a.users[id] = a.userSwitch(user, newUser)
		}
	}
}
