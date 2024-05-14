package usecase

import (
	"music_bot/internal/entity"
)

type HostUser struct {
	user      *entity.HostUser
	sender    Sender
	audioRepo AudioRepo

	updateCh    <-chan entity.Update
	msgCh       <-chan entity.UserMsg
	childMsgChs map[entity.UserID]chan<- entity.UserMsg
}

func (u *HostUser) run(stop chan<- entity.UserID) {
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

func (u *HostUser) handleMsg(update entity.UserMsg) (stop bool) {

	return false
}

func (u *HostUser) handleUpdate(update entity.Update) (stop bool) {
	if update.Command != "" {
		if u.handleCommand(update.Command) {
			return true
		}
	}
	return false
}

func (u *HostUser) handleCommand(command string) (stop bool) {
	switch command {
	case "/menu":
		u.sendMenu()
	case "/start":
		u.out()
		return true
	case "/info":
		u.sendInfo()
	case "/get_tracks":
		u.getTracks()
	case "/get_summary":
		u.getSummary()
	case "/send_notify":
		u.sendNotify()
	}
	return false
}

func (u *HostUser) sendMenu() {
	u.sender.SendMenu(u.user.ID)
}

func (u *HostUser) out() {
	for _, ch := range u.childMsgChs {
		ch <- entity.UserMsg{
			MsgId: entity.UserMsgExit,
			From:  u.user.ID,
		}
	}
	u.sender.Out(u.user.ID)
}

func (u *HostUser) sendInfo() {
	u.sender.SendInfo(u.user.ID)
}

func (u *HostUser) getTracks() {}

func (u *HostUser) getSummary() {}

func (u *HostUser) sendNotify() {}

func (u *HostUser) join(id entity.UserID, msgCh chan<- entity.UserMsg) {

}

func (u *HostUser) disconnectUser(id entity.UserID) {

}

func (u *HostUser) findAudio(url string) (audio entity.Audio) {

	return audio
}

func (u *HostUser) setAudio(audio entity.Audio) {

}
