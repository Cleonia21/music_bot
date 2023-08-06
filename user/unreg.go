package user

import (
	"errors"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
	"github.com/withmandala/go-log"
	"strings"
)

type unregUser struct {
	userFather
	url     string
	pass    string
	blocker bool
}

func (u *unregUser) Init(tg Bot, logger *log.Logger, id telego.ChatID) {
	u.fatherInit(tg, logger, id)
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
			telegoutil.InlineKeyboardButton("принимать").WithCallbackData("host"),
			telegoutil.InlineKeyboardButton("отправлять").WithCallbackData("send"),
		),
	)
	text := "👥Ты можешь выбрать одну из *ролей:*\n\n" +
		"👤*Принимать:* твои друзья будут присылать треки, а я буду ставить их в очередь\\. " +
		"Когда ты попросишь я пришлю тебе пакет из треков, по одному от каждого друга\\. " +
		"Таким образом вы сможете слушать общий плейлист\\. " +
		"Есть одно \"Но\", ты не сможешь добавлять треки в общую очередь\\.😏\n\n" +
		"👤*Отправлять:* ты сможешь отправлять треки, они попадут в общую очередь, " +
		"ты услышишь и свои треки, и треки друзей\\."
	msg := telegoutil.Message(
		u.id,
		text,
	).WithReplyMarkup(keyboard)
	u.sendMessage(msg, true)
	u.clearData()
}

func (u *unregUser) handler(update *telego.Update) (user users, needInit bool) {
	if update.CallbackQuery != nil {
		u.blocker = false

		_, err := u.tg.EditMessageText(&telego.EditMessageTextParams{
			ChatID:    telego.ChatID{ID: update.CallbackQuery.Message.Chat.ID},
			MessageID: update.CallbackQuery.Message.MessageID,
			Text:      update.CallbackQuery.Message.Text,
		})
		if err != nil {
			u.logger.Errorf(err.Error())
		}

		if update.CallbackQuery.Data == "host" {
			hUser := hostUser{}
			return &hUser, true
		} else if update.CallbackQuery.Data == "send" {
			u.sendText("Пришли secretMessage", false)
			u.clearData()
		} else {
			u.logger.Errorf("data not found")
			text := "Неизвестная ошибка на стороне сервера,\nпопробуй нажать /start"
			u.sendText(text, false)
		}

	} else if update.Message != nil {
		if update.Message.Text == "/menu" || update.Message.Text == "/start" {
			u.sendFirstMenu()
			return
		}
		if u.blocker {
			u.sendFirstMenu()
			return
		}

		if err := u.parseSecretMsg(update.Message.Text); err != nil {
			u.sendText(err.Error(), false)
			u.sendFirstMenu()
		} else {
			return &sendingUser{}, true
		}
	}
	return nil, false
}

func (u *unregUser) parseSecretMsg(text string) error {
	err := errors.New("не верный формат секретного сообщения")
	strs := strings.Split(text, "/")
	if len(strs) != 3 {
		return err
	}
	if strs[0] != "secretMessage" {
		return err
	}
	u.url = strs[1][1:]
	u.pass = strs[2]
	return nil
}

func (u *unregUser) notValidate() {
	u.sendText("ссылка или пароль не верные", false)
	u.sendFirstMenu()
}
