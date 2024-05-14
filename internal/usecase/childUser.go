package usecase

import "music_bot/internal/entity"

type ChildUser struct {
	user      *entity.ChildUser
	sender    Sender
	audioRepo AudioRepo

	updateCh   <-chan entity.Update
	msgCh      <-chan entity.UserMsg
	hostMsgChs chan<- entity.UserMsg
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
