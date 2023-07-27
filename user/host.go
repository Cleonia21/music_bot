package user

import (
	"MusicBot/user/playList"
	"fmt"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
	"github.com/withmandala/go-log"
)

type hostUser struct {
	userFather
	pass          string
	connectedUser map[*sendingUser]struct{}
	playList      playList.PlayList
}

func (h *hostUser) init(tg *telego.Bot, logger *log.Logger, chatID telego.ChatID) {
	h.fatherInit(tg, logger, chatID)

	h.pass = "test pass"
	h.connectedUser = make(map[*sendingUser]struct{})
	h.playList.Init()
	text := "Отправь секретное сообщение тем кто хочет присоедениться⤵️"
	h.sendMessage(telegoutil.Message(h.id, text))
	h.sendText(fmt.Sprintf("`secretMessage/@%v/%v`", h.id.Username, h.pass))
}

func (h *hostUser) handler(update *telego.Update) (user users, needInit bool) {
	if update.Message != nil {
		switch update.Message.Text {
		case "/menu":
			h.sendMenu()
		case "/start": //сообщить полям
			return &unregUser{}, true
		}
	}
	if update.CallbackQuery != nil {
		if update.CallbackQuery.Data == "getTracks" {
			h.sendAudioPack()
			_ = h.tg.AnswerCallbackQuery(
				&telego.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})
		}
	}
	return h, false
}

func (h *hostUser) sendMenu() {
	text := "Количество треков у пользователей:\n\n"
	summary := h.playList.GetSummary()
	for _, s := range summary {
		text += fmt.Sprintf("%v\\(%v\\)\n", s.ID.Username, s.Num)
	}
	keyboard := telegoutil.InlineKeyboard(
		telegoutil.InlineKeyboardRow(
			telegoutil.InlineKeyboardButton("дай мне порцию треков").WithCallbackData("getTracks"),
		),
	)
	h.sendMessage(telegoutil.Message(h.id, text).WithReplyMarkup(keyboard))
}

func (h *hostUser) sendAudioPack() (sentMsgs []*telego.Message) {
	audios, _ := h.playList.GetAudio()
	for _, audio := range audios {
		msg := h.sendAudio(audio)
		if msg != nil {
			sentMsgs = append(sentMsgs, msg)
		}
	}
	return sentMsgs
}

func (h *hostUser) validatePass(pass string) (ok bool) {
	if pass == h.pass {
		return true
	} else {
		return false
	}
}

func (h *hostUser) join(user *sendingUser) (sentMsg *telego.Message) {
	h.connectedUser[user] = struct{}{}
	sentMsg = h.sendText("Присоединился @" + user.id.Username)
	return
}

func (h *hostUser) disconnectUser(user *sendingUser) (sentMsg *telego.Message) {
	delete(h.connectedUser, user)
	sentMsg = h.sendText("Пользователь @" + user.id.Username + " отключился")
	return
}

func (h *hostUser) setAudio(from *sendingUser, audio *telego.SendAudioParams) (sentMsg *telego.Message, err error) {
	audio.ChatID = h.id
	err = h.playList.SetAudio(from.id, audio)
	if err != nil {
		return
	}
	sentMsg = h.sendText("@" + from.id.Username + " добавил в очередь трек " + audio.Title)
	return
}
