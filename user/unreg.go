package user

import (
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
)

type unregUser struct {
	userFather
	url     string
	pass    string
	blocker bool
}

func (u *unregUser) Init(tg *telego.Bot, id telego.ChatID) {
	u.tg = tg
	u.id = id
	u.sendFirstMenu()
}

func (u *unregUser) clearData() {
	u.url = ""
	u.pass = ""
}

func (u *unregUser) sendFirstMenu() {
	u.blocker = true
	keyboard := telegoutil.InlineKeyboard(
		telegoutil.InlineKeyboardRow(
			telegoutil.InlineKeyboardButton("принимать треки").WithCallbackData("host"),
		),
		telegoutil.InlineKeyboardRow(
			telegoutil.InlineKeyboardButton("отправлять треки").WithCallbackData("send"),
		),
	)
	msg := telegoutil.Message(
		u.id,
		"выбери вариант",
	).WithReplyMarkup(keyboard)
	u.sendMessage(msg)
	u.clearData()
}

func (u *unregUser) handler(update *telego.Update) (user users, needInit bool) {
	if update.CallbackQuery != nil {
		u.blocker = false
		_ = u.tg.AnswerCallbackQuery(
			&telego.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})

		if update.CallbackQuery.Data == "host" {
			hUser := hostUser{}
			return &hUser, true
		}
		if update.CallbackQuery.Data == "send" {
			u.sendText("пришли ссылку")
			u.clearData()
		}
	} else if update.Message != nil {
		if update.Message.Text == "/menu" || update.Message.Text == "/start" {
			u.sendFirstMenu()
			return
		}
		if u.blocker {
			u.sendText("Я жду нажатие на кнопку или встроенную команду. Не ломай меня, пожалуйста)")
			u.sendFirstMenu()
			return
		}
		if u.url == "" {
			u.url = update.Message.Text[1:]
			u.sendText("пришли пароль")
		} else if u.url != "" {
			u.pass = update.Message.Text
			sUser := sendingUser{}
			return &sUser, true
		}
	}
	return nil, false
}

func (u *unregUser) notValidate() {
	u.sendText("ссылка или пароль не верные")
	u.sendFirstMenu()
}
