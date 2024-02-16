package user

import (
	"MusicBot/log"
	"MusicBot/telegram"
	"MusicBot/user/utils"
	utils2 "MusicBot/utils"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
)

var testFlag = false

type userFather struct {
	id          utils.UserID
	getFromTgCh <-chan telego.Update
}

func (u *userFather) fatherInit(id utils.UserID) {
	u.id = id
}

func (u *userFather) sendText(text string, notification bool) (sentMsg *telego.Message) {
	if text == "" {
		return
	}
	msg := telegoutil.Message(u.id.ChatID, text)

	if testFlag {
		log.Logger.Debugf(utils2.MsgParamsToStr(msg))
		return nil
	}

	return u.sendMessage(msg, notification)
}

func (u *userFather) sendMessage(msg *telego.SendMessageParams, notification bool) (sentMsg *telego.Message) {
	if testFlag {
		log.Logger.Debugf(utils2.MsgParamsToStr(msg))
		return nil
	}

	if !notification {
		msg.WithDisableNotification()
	}
	msg.WithParseMode("HTML")
	sentMsg, err := telegram.TG.SendMessage(msg)
	if err != nil {
		log.Logger.Errorf(err.Error())
	} else {
		log.Logger.Debugf(utils2.MsgToStr(sentMsg))
	}
	return
}

func (u *userFather) sendAudioToUser(audio *telego.SendAudioParams) (sentMsg *telego.Message) {
	if testFlag {
		log.Logger.Debugf(utils2.AudioParamsToStr(audio))
		return nil
	}

	audio.WithDisableNotification()
	msg, err := telegram.TG.SendAudio(audio)
	if err != nil {
		log.Logger.Errorf(err.Error())
	} else {
		log.Logger.Debugf(utils2.MsgToStr(msg))
	}
	return msg
}

func (u *userFather) getID() utils.UserID {
	return u.id
}

func (u *userFather) sendInfo() {
	text := "Ты можешь:\n" +
		"❕Отправить ссылку на трек из Яндекс Музыки\n" +
		"❕Прислать аудиофайл со своего устройства\n" +
		"❕Найти трек в дружественном <a href='https://t.me/YaMuBbot'>боте</a> и переслать мне(потребуется авторизация)"
	msg := telegoutil.Message(u.id.ChatID, text)

	if testFlag {
		log.Logger.Debugf(utils2.MsgParamsToStr(msg))
	} else {
		u.sendMessage(msg, true)
	}
}
