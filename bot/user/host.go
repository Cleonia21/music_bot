package user

import (
	Audio "MusicBot/audio"
	"MusicBot/log"
	"MusicBot/passGen"
	"MusicBot/telegram"
	"MusicBot/user/playList"
	"MusicBot/user/utils"
	utils2 "MusicBot/utils"
	"fmt"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
	"sync"
)

type msgBetweenUsers struct {
	id   string
	from utils.UserID

	text  string
	audio telego.Audio
	url   string
}

type hostUser struct {
	userFather
	pass           string
	playList       playList.PlayList
	audio          *Audio.Audio
	getFromUsersCh chan msgBetweenUsers
	usersMapMutex  sync.RWMutex
	sendToUsersChs map[utils.UserID]chan<- msgBetweenUsers
}

func (h *hostUser) init(chatID utils.UserID, audio *Audio.Audio, getFromTgCh <-chan telego.Update) {
	h.id = chatID
	h.pass = passGen.GeneratePassword(10, 3, 2, 2)
	h.playList.Init()
	h.audio = audio
	h.getFromUsersCh = make(chan msgBetweenUsers, 10)
	h.sendToUsersChs = make(map[utils.UserID]chan<- msgBetweenUsers)
	h.getFromTgCh = getFromTgCh

	h.sendText("Вернуться в начало: /start\nУправление ботом: /menu\nКак прислать музыку: /info", false)
	h.sendText("Ты принимаешь треки👍", false)
	h.sendText("Отправь секретное сообщение тем кто хочет присоедениться⤵️", false)
	h.sendText(fmt.Sprintf("<code>secretMessage/@%v/%v</code>", h.id.ChatID.Username, h.pass), false)
}

func (h *hostUser) join(id utils.UserID, senderInCh chan<- msgBetweenUsers) (hostInCh chan<- msgBetweenUsers) {
	h.usersMapMutex.Lock()
	h.sendToUsersChs[id] = senderInCh
	h.usersMapMutex.Unlock()

	h.sendText(utils.UserNameInserting("Присоединился ", id, ""), false)
	return h.getFromUsersCh
}

func (h *hostUser) out() {
	for _, ch := range h.sendToUsersChs {
		ch <- msgBetweenUsers{id: "out", from: h.id}
	}
	h.sendText("Ты вышел из роли", false)
}

func (h *hostUser) disconnectUser(id utils.UserID) {
	delete(h.sendToUsersChs, id)
	h.sendText(utils.UserNameInserting("Пользователь ", id, " отключился"), false)
}

func (h *hostUser) handler(update *telego.Update) (user users, needInit bool) {
	if update.Message != nil {
		switch update.Message.Text {
		case "/menu":
			h.sendMenu()
		case "/start":
			h.out()
			return &unregUser{}, true
		case "/info":
			h.sendInfo()
		default:
			h.setAudioToPlaylist(update)
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
		_ = telegram.TG.AnswerCallbackQuery(
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
	h.sendMessage(telegoutil.Message(h.id.ChatID, text).WithReplyMarkup(keyboard), true)
}

func (h *hostUser) setAudioToPlaylist(update *telego.Update) {
	track, err := h.audio.GetParams(update)
	if err != nil {
		h.sendText("Не удалось получить трек", false)
		log.Logger.Errorf("err: %v, update: %v", err.Error(), utils2.UpdateToStr(update))
	} else {
		track.ChatID = h.id.ChatID
		err = h.playList.SetAudio(h.id, track)
		if err != nil {
			h.sendText(err.Error(), false)
		} else {
			h.sendText("Отправил в очередь", false)
		}
	}
}

func (h *hostUser) setAudioToPlaylistFromUser(id utils.UserID, audio *telego.SendAudioParams) (sentMsg *telego.Message, err error) {
	audio.ChatID = h.id.ChatID
	err = h.playList.SetAudio(id, audio)
	if err != nil {
		return
	}
	sentMsg = h.sendText(utils.UserNameInserting(
		"",
		id,
		" добавил в очередь трек: \""+audio.Title+"\""),
		false)
	return
}

func (h *hostUser) sendAudioPack() (sentMsgs []*telego.Message) {
	audios, _ := h.playList.GetAudio()
	for _, audio := range audios {
		msg := h.sendAudioToUser(audio)
		if msg != nil {
			sentMsgs = append(sentMsgs, msg)
		}
	}
	return sentMsgs
}

func (h *hostUser) sendNotify() (sentMsg *telego.Message) {
	summary := h.playList.GetSummary()
	if len(summary) == 0 || (len(summary) == 1 && summary[0].ID == h.id) {
		return h.sendText("Никто еще не подключился", false)
	}
	for _, sum := range summary {
		if sum.Num < 2 && sum.ID != h.id {
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

func (h *hostUser) trackNum(who utils.UserID) int {
	return h.playList.UserTrackNum(who)
}

func (h *hostUser) validatePass(pass string) (ok bool) {
	if pass == h.pass {
		return true
	} else {
		return false
	}
}
