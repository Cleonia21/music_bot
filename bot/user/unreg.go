package user

import (
	"MusicBot/log"
	"MusicBot/telegram"
	"MusicBot/user/utils"
	"errors"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
	"strings"
)

type unregUser struct {
	userFather
	url     string
	pass    string
	blocker bool
}

func (u *unregUser) Init(id utils.UserID) {
	u.fatherInit(id)
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
	text := "👥Ты можешь выбрать одну из <b>ролей:</b>\n\n" +
		"👤<b>Принимать:</b> ты и твои друзья будут присылать треки, а я буду ставить их в очередь. " +
		"Когда ты попросишь я пришлю тебе пакет из треков, по одному от каждого друга. " +
		"Таким образом вы сможете слушать общий плейлист.\n\n" +
		"👤<b>Отправлять:</b> ты сможешь отправлять треки, они попадут в общую очередь, " +
		"ты услышишь и свои треки, и треки друзей."
	msg := telegoutil.Message(
		u.id.ChatID,
		text,
	).WithReplyMarkup(keyboard)
	u.sendMessage(msg, false)
	u.clearData()
}

func (u *unregUser) handler(update *telego.Update) (user users, needInit bool) {
	if update.CallbackQuery != nil {
		u.blocker = false

		_, err := telegram.TG.EditMessageText(&telego.EditMessageTextParams{
			ChatID:    telego.ChatID{ID: update.CallbackQuery.Message.Chat.ID},
			MessageID: update.CallbackQuery.Message.MessageID,
			Text:      update.CallbackQuery.Message.Text,
		})
		if err != nil {
			log.Logger.Errorf(err.Error())
		}

		if update.CallbackQuery.Data == "host" {
			hUser := hostUser{}
			return &hUser, true
		} else if update.CallbackQuery.Data == "send" {
			u.sendText("Пришли secretMessage", false)
			u.clearData()
		} else {
			log.Logger.Errorf("data not found")
			text := "Неизвестная ошибка на стороне сервера,\nпопробуй нажать /start"
			u.sendText(text, false)
		}

	} else if update.Message != nil {
		if update.Message.Text == "/start" {
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
	u.sendText("не верное secretMessage", false)
	u.sendFirstMenu()
}
