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

type Admin struct {
	tg     *telego.Bot
	users  map[int64]users
	audio  *audio.Audio
	logger *log.Logger
}

func Init(tg *telego.Bot) *Admin {
	l := log.New(os.Stderr)
	l.WithColor()
	l.WithDebug()

	a := Admin{
		tg:     tg,
		users:  make(map[int64]users),
		audio:  audio.Init(),
		logger: l,
	}
	return &a
}

func (a *Admin) Handler(update *telego.Update) {
	a.logger.Debugf("get update: " + utils.UpdateToStr(update))

	var from int64
	if update.Message != nil {
		from = update.Message.From.ID
	} else if update.CallbackQuery != nil {
		from = update.CallbackQuery.From.ID
	}

	user, ok := a.users[from]
	if !ok {
		user := &unregUser{}
		user.Init(
			a.tg,
			a.logger,
			telego.ChatID{
				ID:       update.Message.Chat.ID,
				Username: update.Message.From.Username,
			},
		)
		a.users[from] = users(user)
	} else {
		newUser, needInit := user.handler(update)
		if needInit {
			a.users[from] = a.init(user, newUser)
		}
	}
}

func (a *Admin) init(from, to users) users {
	var unreg *unregUser
	var sending *sendingUser
	var host *hostUser
	ok := false

	if unreg, ok = from.(*unregUser); ok {
		if host, ok = to.(*hostUser); ok {
			host.init(
				a.tg,
				a.logger,
				unreg.id,
			)
			return host
		} else if sending, ok = to.(*sendingUser); ok {
			host, ok = a.searchUser(telego.ChatID{Username: unreg.url}).(*hostUser)
			if ok && host.validatePass(unreg.pass) {
				sending.init(
					a.tg,
					a.logger,
					unreg.id,
					host,
					a.audio,
				)
				host.join(sending)
				return sending
			} else {
				unreg.notValidate()
				return unreg
			}
		}
	} else if unreg, ok = to.(*unregUser); ok {
		if host, ok = from.(*hostUser); ok {
			unreg.Init(
				a.tg,
				a.logger,
				host.id,
			)
		} else if sending, ok = from.(*sendingUser); ok {
			unreg.Init(
				a.tg,
				a.logger,
				sending.id,
			)
		}
		return unreg
	}
	return nil
}

func (a *Admin) searchUser(id telego.ChatID) users {
	for _, user := range a.users {
		if id.ID == user.getID().ID ||
			id.Username == user.getID().Username {
			return user
		}
	}
	return nil
}
