package user

import (
	"MusicBot/user/playList"
	"fmt"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
)

type hostUser struct {
	userFather
	pass          string
	connectedUser map[*sendingUser]struct{}
	playList      playList.PlayList
}

func (h *hostUser) init(tg *telego.Bot, chatID telego.ChatID) {
	h.tg = tg
	h.id = chatID
	h.pass = "test pass"
	h.connectedUser = make(map[*sendingUser]struct{})
	h.playList.Init()
	h.sendText(fmt.Sprintf("Имя: @%v\nПароль: %v", h.id.Username, h.pass))
	h.sendMessage(telegoutil.Message(h.id, "Ждем пока кто-нибудь присоеденится..."))
}

func (h *hostUser) handler(update *telego.Update) (user users, needInit bool) {
	if update.Message != nil {
		switch update.Message.Text {
		case "/menu":
			h.sendMenu()
		case "/start":
			return &unregUser{}, true
		}
	}
	if update.CallbackQuery != nil {
		_ = h.tg.AnswerCallbackQuery(
			&telego.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})

		if update.CallbackQuery.Data == "getTracks" {
			h.sendAudio()
		}
	}
	return h, false
}

func (h *hostUser) sendMenu() {
	text := "Количество треков у пользователей:\n\n"
	summary := h.playList.GetSummary()
	for _, s := range summary {
		text += fmt.Sprintf("%v(%v)\n", s.ID.Username, s.Num)
	}
	keyboard := telegoutil.InlineKeyboard(
		telegoutil.InlineKeyboardRow(
			telegoutil.InlineKeyboardButton("дай мне порцию треков").WithCallbackData("getTracks"),
		),
	)
	h.sendMessage(telegoutil.Message(h.id, text).WithReplyMarkup(keyboard))
}

func (h *hostUser) sendAudio() {
	audios, errs := h.playList.GetAudio()
	text := ""
	if len(errs) != 0 {
		text += "пользователи не добавили треки:\n"
	}
	for _, err := range errs {
		text += err + " "
	}
	h.sendText(text)
	for _, audio := range audios {
		_, err := h.tg.SendAudio(audio)
		if err != nil {
			h.tg.Logger().Errorf(err.Error())
		}
	}
}

func (h *hostUser) validatePass(pass string) (ok bool) {
	if pass == h.pass {
		return true
	} else {
		return false
	}
}

func (h *hostUser) join(user *sendingUser) {
	h.connectedUser[user] = struct{}{}
	h.sendText("Присоединился @" + user.id.Username)
}

func (h *hostUser) disconnectUser(user *sendingUser) {
	delete(h.connectedUser, user)
	h.sendText("Пользователь @" + user.id.Username + " отключился")
}

func (h *hostUser) setAudio(from *sendingUser, audio *telego.SendAudioParams) error {
	audio.ChatID = h.id
	err := h.playList.SetAudio(from.id, audio)
	if err != nil {
		return err
	}
	h.sendText("@" + from.id.Username + " добавил в очередь трек " + audio.Title)
	return nil
}
