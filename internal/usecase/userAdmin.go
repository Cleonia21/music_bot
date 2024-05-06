package usecase

import (
	"music_bot/internal/entity"
	"sync"
)

type Admin struct {
	users         map[entity.UserID]chan entity.Update
	usersMapMutex sync.RWMutex
	stop          chan entity.UserID
}

func Init() *Admin {
	a := Admin{
		users: make(map[entity.UserID]chan entity.Update),
		stop:  make(chan entity.UserID, 10),
	}
	go a.stopHandler()
	return &a
}

func (a *Admin) stopHandler() {
	for {
		id := <-a.stop
		a.usersMapMutex.Lock()
		delete(a.users, id)
		a.usersMapMutex.Unlock()
	}
}

func (a *Admin) Handler(update entity.Update) {
	updateCh, ok := a.users[update.ID]

	if !ok {
		user := NewUnregUser(update.ID)
		a.users[update.ID] = user.updateCh
	} else {
		updateCh <- update
	}
}
