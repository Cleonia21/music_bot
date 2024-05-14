package usecase

import "music_bot/internal/entity"

type ChildUser struct {
	user      *entity.HostUser
	audioRepo AudioRepo

	updateCh   <-chan entity.Update
	msgCh      chan entity.UserMsg
	hostMsgChs chan<- entity.UserMsg
}

func newChildUser(id entity.UserID, updateCH <-chan entity.Update, hostMsgChs chan<- entity.UserMsg) (u *ChildUser) {
	return nil
}

func (u *ChildUser) run(stop chan<- entity.UserID) {
	// select {
	// case update := <-u.updateCh:
	// 	if !u.handleUpdate(update) {
	// 		stop <- u.user.ID
	// 		return
	// 	}
	// case msg := <-u.msgCh:
	// 	if !u.handleMsg(msg) {
	// 		stop <- u.user.ID
	// 		return
	// 	}
	// }
}