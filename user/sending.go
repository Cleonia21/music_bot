package user

import (
	"MusicBot/audio"
	"github.com/mymmrac/telego"
)

type sendingUser struct {
	userFather
	host  *hostUser
	audio *audio.Audio
}

func (s *sendingUser) init(tg *telego.Bot, chatID telego.ChatID, host *hostUser,
	audio *audio.Audio) {

	s.tg = tg
	s.id = chatID
	s.host = host
	s.audio = audio

	s.sendText("Присылай ссылки с яндекс музыки")
}

func (s *sendingUser) connect(user *hostUser) {
	s.host = user
}

func (s *sendingUser) handler(update *telego.Update) (user users, needInit bool) {
	if update.Message != nil {
		switch update.Message.Text {
		case "/start":
			return &unregUser{}, true
		default:
			track, err := s.audio.GetParams(update)
			if err != nil {
				s.sendText("Не удалось получить трек")
			} else {
				err = s.host.setAudio(s, track)
				if err != nil {
					s.sendText(err.Error())
				} else {
					s.sendText("Отправил в очередь")
				}
			}
		}
	}
	return s, false
}
