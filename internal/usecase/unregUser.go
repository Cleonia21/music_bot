package usecase

import "music_bot/internal/entity"

type UnregUser struct {
	updateCh chan entity.Update
}

func NewUnregUser(id entity.UserID) (u *UnregUser) {

	return u
}

func (u *UnregUser) Run(stop chan<- entity.UserID) {

}
