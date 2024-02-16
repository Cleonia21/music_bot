package user

import "MusicBot/user/utils"

type post struct {
	// users хранит пользователей, которым они отправляют сообщения
	users map[utils.UserID]map[utils.UserID]interface{}
}

func (p *post) post() {

}

func (p *post) connect(who utils.UserID, whom utils.UserID) (<-chan msgBetweenUsers, chan<- msgBetweenUsers) {

}
