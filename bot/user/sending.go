package user

import (
	"MusicBot/audio"
	"MusicBot/log"
	"MusicBot/user/utils"
	utils2 "MusicBot/utils"
	"fmt"
	"github.com/mymmrac/telego"
)

type sendingUser struct {
	userFather
	host  *hostUser
	audio *audio.Audio
}

func (s *sendingUser) init(chatID utils.UserID, host *hostUser,
	audio *audio.Audio) {

	s.fatherInit(chatID)

	s.host = host
	s.audio = audio

	s.sendText("Ты успешно присоеденился👍", false)
	s.sendText("Вернуться в начало: /start\nУправление ботом: /menu\nКак прислать музыку: /info", false)
}

func (s *sendingUser) connect(user *hostUser) {
	s.host = user
}

func (s *sendingUser) disconnect() {
	s.sendText("Ты вышел из роли", false)
	s.host.disconnectUser(s)
}

func (s *sendingUser) handler(update *telego.Update) (user users, needInit bool) {
	if update.Message != nil {
		if s.host == nil {
			return &unregUser{}, true
		}

		switch update.Message.Text {
		case "/start":
			s.disconnect()
			return &unregUser{}, true
		case "/menu":
			s.sendMenu()
		case "/info":
			s.sendInfo()
		default:
			s.setAudio(update)
		}
	}
	return s, false
}

func (s *sendingUser) sendMenu() {
	s.sendText(fmt.Sprintf("У тебя в очереди еще %v трек(а/ов)", s.host.trackNum(s.id)), false)
}

func (s *sendingUser) setAudio(update *telego.Update) {
	track, err := s.audio.GetParams(update)
	if err != nil {
		s.sendText("Не удалось получить трек", false)
		log.Logger.Errorf("err: %v, update: %v", err.Error(), utils2.UpdateToStr(update))
	} else {
		_, err = s.host.setAudioToPlaylistFromUser(s.id, track)
		if err != nil {
			s.sendText(err.Error(), false)
		} else {
			s.sendText("Отправил в очередь", false)
		}
	}
}

func (s *sendingUser) hostOut() {
	s.host = nil
	s.sendText("Принимающий пользователь отключился, жми /start", false)
}

func (s *sendingUser) tracksEndedInQueue() {
	s.sendText("Принимающий пользователь поросит прислать еще треков", true)
}
