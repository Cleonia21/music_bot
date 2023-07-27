package user

import (
	"MusicBot/user/utils"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
	"github.com/withmandala/go-log"
)

type userFather struct {
	tg     *telego.Bot
	id     telego.ChatID
	logger *log.Logger
}

func (u *userFather) fatherInit(tg *telego.Bot, logger *log.Logger, id telego.ChatID) {
	u.tg = tg
	u.id = id
	u.logger = logger
}

func (u *userFather) sendText(text string) (sentMsg *telego.Message) {
	if text == "" {
		return
	}
	msg := telegoutil.Message(u.id, text).WithParseMode("MarkdownV2").WithDisableNotification()
	sentMsg, err := u.tg.SendMessage(msg)

	if err != nil {
		u.logger.Errorf(err.Error())
	} else {
		u.logger.Debugf(utils.MsgToStr(sentMsg))
	}
	return
}

func (u *userFather) sendMessage(msg *telego.SendMessageParams) (sentMsg *telego.Message) {
	msg.WithParseMode("MarkdownV2").WithDisableNotification()
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
