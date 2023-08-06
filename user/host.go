package user

import (
	"MusicBot/user/playList"
	"MusicBot/user/utils"
	"fmt"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
	"github.com/withmandala/go-log"
)

type hostUser struct {
	userFather
	pass          string
	connectedUser map[utils.UserID]*sendingUser
	playList      playList.PlayList
}

func (h *hostUser) init(tg Bot, logger *log.Logger, chatID utils.UserID) {
	h.fatherInit(tg, logger, chatID)

	h.pass = "test pass"
	h.connectedUser = make(map[utils.UserID]*sendingUser)
	h.playList.Init()
	h.sendText("Ты принимаешь треки👍 Если что, есть команда /menu", false)
	h.sendText("Отправь секретное сообщение тем кто хочет присоедениться⤵️", false)
	h.sendText(fmt.Sprintf("<code>secretMessage/@%v/%v</code>", h.id.ChatID.Username, h.pass), false)
}

func (h *hostUser) handler(update *telego.Update) (user users, needInit bool) {
	if update.Message != nil {
		switch update.Message.Text {
		case "/menu":
			h.sendMenu()
		case "/start":
			h.out()
			return &unregUser{}, true
		}
	}
	if update.CallbackQuery != nil {
		switch update.CallbackQuery.Data {
		case "getTracks":
			h.sendAudioPack()
		case "getSummary":
			h.sendSummary()
		case "sendNotify":
			h.sendNotify()
		}
		_ = h.tg.AnswerCallbackQuery(
			&telego.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})
	}
	return h, false
}

func (h *hostUser) sendMenu() {
	text := "Меню🎛"
	keyboard := telegoutil.InlineKeyboard(
		telegoutil.InlineKeyboardRow(
			telegoutil.InlineKeyboardButton("новая порция 🎧").WithCallbackData("getTracks"),
		),
		telegoutil.InlineKeyboardRow(
			telegoutil.InlineKeyboardButton("количество 🎧 у 👤").WithCallbackData("getSummary"),
		),
		telegoutil.InlineKeyboardRow(
			telegoutil.InlineKeyboardButton("оповестить 👤 у которых мало 🎧").WithCallbackData("sendNotify"),
		),
	)
	h.sendMessage(telegoutil.Message(h.id.ChatID, text).WithReplyMarkup(keyboard), false)
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

func (h *hostUser) sendNotify() (sentMsg *telego.Message) {
	summary := h.playList.GetSummary()
	if len(summary) == 0 {
		return h.sendText("Никто еще не подключился", false)
	}
	for _, sum := range summary {
		if sum.Num < 2 {
			h.connectedUser[sum.ID].tracksEndedInQueue()
		}
	}
	return h.sendText("Отправил уведомление тем у кого 1 или 0 треков", false)
}

func (h *hostUser) sendSummary() (sentMsg *telego.Message) {
	text := "Количество треков у пользователей:\n\n"
	summary := h.playList.GetSummary()
	if len(summary) == 0 {
		return h.sendText("Никто еще ничего не добавил", false)
	}
	for _, s := range summary {
		text += utils.UserNameInserting("", s.ID, fmt.Sprintf("(%v)", s.Num))
	}
	return h.sendText(text, false)
}

func (h *hostUser) validatePass(pass string) (ok bool) {
	if pass == h.pass {
		return true
	} else {
		return false
	}
}

func (h *hostUser) join(user *sendingUser) (sentMsg *telego.Message) {
	h.connectedUser[user.id] = user
	sentMsg = h.sendText(utils.UserNameInserting("Присоединился ", user.id, ""), false)
	return
}

func (h *hostUser) disconnectUser(user *sendingUser) (sentMsg *telego.Message) {
	delete(h.connectedUser, user.id)
	sentMsg = h.sendText(utils.UserNameInserting("Пользователь ", user.id, " отключился"), false)
	return
}

func (h *hostUser) setAudio(from *sendingUser, audio *telego.SendAudioParams) (sentMsg *telego.Message, err error) {
	audio.ChatID = h.id.ChatID
	err = h.playList.SetAudio(from.id, audio)
	if err != nil {
		return
	}
	sentMsg = h.sendText(utils.UserNameInserting(
		"",
		from.id,
		" добавил в очередь трек: \""+audio.Title+"\""),
		false)
	return
}

func (h *hostUser) out() {
	for _, user := range h.connectedUser {
		user.hostOut()
	}
	h.sendText("Ты вышел из роли", false)
}

func (h *hostUser) trackNum(who utils.UserID) int {
	return h.playList.UserTrackNum(who)
}
