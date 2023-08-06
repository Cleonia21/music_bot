package user

import (
	"MusicBot/user/utils"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
	"github.com/withmandala/go-log"
)

type userID struct {
	telego.ChatID
	firstName string
}

type userFather struct {
	tg     Bot
	id     telego.ChatID
	logger *log.Logger
}

func (u *userFather) fatherInit(tg Bot, logger *log.Logger, id telego.ChatID) {
	u.tg = tg
	u.id = id
	u.logger = logger
}

func (u *userFather) sendText(text string, mode bool) (sentMsg *telego.Message) {
	if text == "" {
		return
	}
	msg := telegoutil.Message(u.id, text)
	return u.sendMessage(msg, mode)
}

/*
msg := telegoutil.Message(telegoutil.ID(update.Message.From.ID), fmt.Sprintf("[inline mention of a user](tg://user?id=%v)", update.Message.From.ID))
msg.WithParseMode("MarkdownV2")
_, _ = tg.SendMessage(msg)
*/

func (u *userFather) sendMessage(msg *telego.SendMessageParams, mode bool) (sentMsg *telego.Message) {
	msg.WithDisableNotification()
	if mode {
		msg.WithParseMode("MarkdownV2")
	}
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

func (u *userFather) getID() telego.ChatID {
	return u.id
}
