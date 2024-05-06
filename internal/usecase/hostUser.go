package usecase

import (
	"music_bot/internal/entity"
)

type HostUser struct {
	user      *entity.HostUser
	audioRepo AudioRepo

	updateCh    <-chan entity.Update
	msgCh       chan entity.UserMsg
	childMsgChs map[entity.UserID]chan<- entity.UserMsg
}

func NewHostUser(id entity.UserID, updateCH <-chan entity.Update, audioRepo AudioRepo) (u *HostUser) {

	return u
}

func (u *HostUser) Run(stop chan<- entity.UserID) {
	select {
	case update := <-u.updateCh:
		if !u.handleUpdate(update) {
			stop <- u.user.ID
			return
		}
	case msg := <-u.msgCh:
		if !u.handleMsg(msg) {
			stop <- u.user.ID
			return
		}
	}
}

func (u *HostUser) handleUpdate(update entity.Update) (stop bool) {

	return false
}

func (u *HostUser) handleMsg(update entity.UserMsg) (stop bool) {

	return false
}

func (u *HostUser) getMsgCh() chan<- entity.UserMsg {
	return u.msgCh
}

func (u *HostUser) join(id entity.UserID, msgCh chan<- entity.UserMsg) {

}

func (u *HostUser) out() {

}

func (u *HostUser) disconnectUser(id entity.UserID) {

}

func (u *HostUser) findAudio(url string) (audio entity.Audio) {

	return audio
}

func (u *HostUser) setAudio(audio entity.Audio) {

}
