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
	getID() telego.ChatID
}

type Bot interface {
	SendMessage(params *telego.SendMessageParams) (*telego.Message, error)
	SendAudio(params *telego.SendAudioParams) (*telego.Message, error)
	AnswerCallbackQuery(params *telego.AnswerCallbackQueryParams) error
	EditMessageText(params *telego.EditMessageTextParams) (*telego.Message, error)
}

type Admin struct {
	tg     Bot
	users  map[telego.ChatID]users
	audio  *audio.Audio
	logger *log.Logger
}

func Init(tg Bot) *Admin {
	l := log.New(os.Stderr)
	l.WithColor()
	l.WithDebug()

	a := Admin{
		tg:     tg,
		users:  make(map[telego.ChatID]users),
		audio:  audio.Init(tg, l),
		logger: l,
	}
	return &a
}

func (a *Admin) Handler(update *telego.Update) {
	a.logger.Debugf("get update: " + utils.UpdateToStr(update))

	userID := utils.UpdateToID(update)
	user, ok := a.users[userID]

	if !ok {
		user := &unregUser{}
		user.Init(
			a.tg,
			a.logger,
			telego.ChatID{ID: update.Message.Chat.ID, Username: update.Message.From.Username},
		)
		a.users[userID] = users(user)
	} else {
		newUser, needInit := user.handler(update)
		if needInit {
			a.users[userID] = a.userSwitch(user, newUser)
		}
	}
}
