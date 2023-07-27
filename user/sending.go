package user

import (
	"MusicBot/audio"
	"MusicBot/user/utils"
	"github.com/mymmrac/telego"
	"github.com/withmandala/go-log"
)

type sendingUser struct {
	userFather
	host  *hostUser
	audio *audio.Audio
}

func (s *sendingUser) init(tg *telego.Bot, logger *log.Logger, chatID telego.ChatID, host *hostUser,
	audio *audio.Audio) {

	s.fatherInit(tg, logger, chatID)

	s.host = host
	s.audio = audio

	s.sendText("Присылай ссылки с яндекс музыки")
}

func (s *sendingUser) connect(user *hostUser) {
	s.host = user
}

func (s *sendingUser) sendAudio(update *telego.Update) {
	track, err := s.audio.GetParams(update)
	if err != nil {
		s.sendText("Не удалось получить трек")
		s.logger.Errorf("err: %v, update: %v", err.Error(), utils.UpdateToStr(update))
	} else {
		_, err = s.host.setAudio(s, track)
		if err != nil {
			s.sendText(err.Error())
		} else {
			s.sendText("Отправил в очередь")
		}
	}
}

func (s *sendingUser) handler(update *telego.Update) (user users, needInit bool) {
	if update.Message != nil {
		switch update.Message.Text {
		case "/start":
			return &unregUser{}, true
		case "/menu":
			s.sendText("Эта команда еще не реализованна")
		default:
			s.sendAudio(update)
		}
	}
	return s, false
}
