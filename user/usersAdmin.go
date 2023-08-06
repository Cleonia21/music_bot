package user

import (
	"MusicBot/audio"
	"MusicBot/user/utils"
	"github.com/mymmrac/telego"
	"github.com/withmandala/go-log"
	"os"
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
	tg     Bot
	users  map[utils.UserID]users
	audio  *audio.Audio
	logger *log.Logger
}

func Init(tg Bot) *Admin {
	l := log.New(os.Stderr)
	l.WithColor()
	l.WithDebug()

	a := Admin{
		tg:     tg,
		users:  make(map[utils.UserID]users),
		audio:  audio.Init(tg, l),
		logger: l,
	}
	return &a
}

func (a *Admin) Handler(update *telego.Update) {
	a.logger.Debugf("get update: " + utils.UpdateToStr(update))

	id := utils.UpdateToID(update)
	user, ok := a.users[id]

	if !ok {
		user := &unregUser{}
		user.Init(
			a.tg,
			a.logger,
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
