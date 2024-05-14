package usecase

import (
	"music_bot/internal/entity"
	"sync"
)

type UserData struct {
	updateCh chan entity.Update
	msgCh    chan entity.UserMsg
	pass     string
}

type Admin struct {
	users         map[entity.UserID]UserData
	unregUsers    map[entity.UserID]*UnregUser
	usersMapMutex sync.RWMutex
	stop          chan entity.UserID
}

func NewAdmin() *Admin {
	a := Admin{
		users: make(map[entity.UserID]UserData),
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
	userChs, ok := a.users[update.UserID]

	if !ok {
		a.unreg(update)
	} else {
		userChs.updateCh <- update
	}
}

func (a *Admin) unreg(update entity.Update) {
	unregUser, ok := a.unregUsers[update.UserID]
	if !ok {
		unregUser = newUnregUser(update.UserID)
		a.unregUsers[update.UserID] = unregUser
	} else {
		action, data := unregUser.setUpdate(update)
		switch action {
		case "connect":
			// поменять ключ
			hostUserData := a.users[update.UserID]
			if hostUserData.pass == data {
				newChildUpdateCh := new(chan entity.Update)
				newChildUser := newChildUser(update.UserID, *newChildUpdateCh, hostUserData.msgCh)
				a.users[update.UserID] = UserData{
					updateCh: *newChildUpdateCh,
					msgCh:    newChildUser.msgCh,
				}
				go newChildUser.run(a.stop)
			}
		case "role":
			if data == "host" {
				newHostUpdateCh := new(chan entity.Update)
				newHostUser := newHostUser(update.UserID, *newHostUpdateCh)
				a.users[update.UserID] = UserData{
					updateCh: *newHostUpdateCh,
					msgCh:    newHostUser.msgCh,
					pass:			newHostUser.getPass(),
				}
				go newHostUser.run(a.stop)
			}
		}
	}
}
