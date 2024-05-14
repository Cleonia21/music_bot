package usecase

import "music_bot/internal/entity"

type UnregUser struct {
	updateCh chan entity.Update
}

func newUnregUser(id entity.UserID) (u *UnregUser) {

	return u
}

func (u *UnregUser) setUpdate(update entity.Update) (action, data string) {

	return "", ""
}

func invalidPass() {

}