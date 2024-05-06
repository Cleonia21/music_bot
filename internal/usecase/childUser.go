package usecase

import "music_bot/internal/entity"

type ChildUser struct {
	user      *entity.HostUser
	audioRepo AudioRepo

	updateCh   <-chan entity.Update
	msgCh      chan entity.UserMsg
	hostMsgChs chan<- entity.UserMsg
}
