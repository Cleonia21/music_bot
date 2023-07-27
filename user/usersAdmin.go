package user

import (
	"MusicBot/audio"
	"github.com/mymmrac/telego"
)

type users interface {
	handler(update *telego.Update) (users, bool)
	getID() telego.ChatID
}

type Admin struct {
	tg    *telego.Bot
	users map[int64]users
	audio *audio.Audio
}

func Init(tg *telego.Bot) *Admin {
	a := Admin{
		tg:    tg,
		users: make(map[int64]users),
		audio: audio.Init(),
	}
	return &a
}

func (a *Admin) Handler(update *telego.Update) {
	var from int64
	if update.Message != nil {
		from = update.Message.From.ID
	} else if update.CallbackQuery != nil {
		from = update.CallbackQuery.From.ID
	}

	if update.Message != nil && update.Message.Text == "/start" {
		user := &unregUser{}
		user.Init(
			a.tg,
			telego.ChatID{
				ID:       update.Message.Chat.ID,
				Username: update.Message.From.Username,
			},
		)
		a.users[from] = users(user)
	} else {
		user, ok := a.users[from]
		if ok {
			newUser, needInit := user.handler(update)
			if needInit {
				a.users[from] = a.init(user, newUser)
			}
		} else {
			a.tg.Logger().Errorf("user not found")
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
			host.init(unreg.tg, unreg.id)
			return host
		} else if sending, ok = to.(*sendingUser); ok {
			host, ok = a.searchUser(telego.ChatID{Username: unreg.url}).(*hostUser)
			if ok && host.validatePass(unreg.pass) {
				sending.init(unreg.tg, unreg.id, host, a.audio)
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
				host.id,
			)
		} else if sending, ok = from.(*sendingUser); ok {
			unreg.Init(
				a.tg,
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
