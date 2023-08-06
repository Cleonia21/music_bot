package user

import (
	"MusicBot/user/utils"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
	"github.com/withmandala/go-log"
)

type userFather struct {
	tg     Bot
	id     utils.UserID
	logger *log.Logger
}

func (u *userFather) fatherInit(tg Bot, logger *log.Logger, id utils.UserID) {
	u.tg = tg
	u.id = id
	u.logger = logger
}

func (u *userFather) sendText(text string, notification bool) (sentMsg *telego.Message) {
	if text == "" {
		return
	}
	msg := telegoutil.Message(u.id.ChatID, text)
	return u.sendMessage(msg, notification)
}

func (u *userFather) sendMessage(msg *telego.SendMessageParams, notification bool) (sentMsg *telego.Message) {
	if !notification {
		msg.WithDisableNotification()
	}
	msg.WithParseMode("HTML")
	sentMsg, err := u.tg.SendMessage(msg)
	if err != nil {
		u.logger.Errorf(err.Error())
	} else {
		u.logger.Debugf(utils.MsgToStr(sentMsg))
	}
	return
}

func (u *userFather) sendAudio(audio *telego.SendAudioParams) (sentMsg *telego.Message) {
	audio.WithDisableNotification()
	msg, err := u.tg.SendAudio(audio)
	if err != nil {
		u.logger.Errorf(err.Error())
	} else {
		u.logger.Debugf(utils.MsgToStr(msg))
	}
	return msg
}

func (u *userFather) getID() utils.UserID {
	return u.id
}
