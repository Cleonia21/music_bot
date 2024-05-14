package usecase

import (
	"errors"
	"music_bot/internal/entity"
	"sync"
)

type UserData struct {
	updateCh chan entity.Update
	msgCh    chan entity.UserMsg
	pass     string
}

type Admin struct {
	audioRepo     AudioRepo
	sender        Sender
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
		a.newUnregUser(update.UserID)
	} else {
		action := unregUser.setUpdate(update)
		switch action {
		case "child":
			hostUserData, ok := a.users[unregUser.hostId]
			if !ok {
				unregUser.setAction("host not found")
			}
			if hostUserData.pass == unregUser.pass {
				if err := a.newChildUser(unregUser.user.ID, unregUser.hostId); err != nil {
					unregUser.setAction(err.Error())
				} else {
					delete(a.unregUsers, unregUser.user.ID)
				}
			} else {
				unregUser.setAction("pass incorrect")
			}
		case "host":
			a.newHostUser(unregUser.user.ID)
			delete(a.unregUsers, unregUser.user.ID)
		}
	}
}

func (a *Admin) newUnregUser(userID entity.UserID) {
	user := UnregUser{
		user:   entity.NewUnregUser(userID),
		sender: a.sender,
	}
	a.unregUsers[userID] = &user
}

func (a *Admin) newChildUser(userID, hostID entity.UserID) error {
	updateCh := new(chan entity.Update)
	msgCh := new(chan entity.UserMsg)

	a.usersMapMutex.Lock()
	defer a.usersMapMutex.Unlock()
	host, ok := a.users[hostID]
	if !ok {
		return errors.New("host not found")
	}

	user := ChildUser{
		user:       entity.NewChildUser(userID),
		sender:     a.sender,
		audioRepo:  a.audioRepo,
		updateCh:   *updateCh,
		msgCh:      *msgCh,
		hostMsgChs: host.msgCh,
	}

	a.users[userID] = UserData{
		updateCh: *updateCh,
		msgCh:    *msgCh,
	}
	go user.run(a.stop)
	return nil
}

func (a *Admin) newHostUser(userID entity.UserID) {
	updateCh := new(chan entity.Update)
	msgCh := new(chan entity.UserMsg)

	a.usersMapMutex.Lock()
	defer a.usersMapMutex.Unlock()

	user := HostUser{
		user:        entity.NewHostUser(userID),
		sender:      a.sender,
		audioRepo:   a.audioRepo,
		updateCh:    *updateCh,
		msgCh:       *msgCh,
		childMsgChs: make(map[entity.UserID]chan<- entity.UserMsg),
	}
	go user.run(a.stop)
}
