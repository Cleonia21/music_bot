package users

import (
	"errors"
	"github.com/mymmrac/telego"
)

type User struct {
	ChatID      telego.ChatID
	ID          int64
	Room        *room
	isRoomRoot  bool
	Nik         string
	MessHistory map[int]string
	State       int
	PlayList    playList
}

func (u *User) Constructor(ChatID telego.ChatID, ID int64, nik string) {
	u.ChatID = ChatID
	u.ID = ID
	u.Nik = nik

	u.MessHistory = make(map[int]string)
	u.PlayList.constructor()
}

func (u *User) Destructor() {
	u.OutRoom()
}

func (u *User) OutRoom() {
	u.Room.out(u)
}

func (u *User) CreateRoom(title, pass string) {
	u.Room = &room{}
	u.Room.constructor(
		u,
		title,
		pass,
	)
	u.isRoomRoot = true
}

func (u *User) AddInRoom(user *User, pass string) error {
	if u.Room == nil {
		return errors.New("the root user has no rooms")
	}
	if u.isRoomRoot == false {
		return errors.New("the \"root\" user cannot control the room")
	}
	if user.Room != nil {
		return errors.New("the user is already in the room")
	}

	err := u.Room.add(user, pass)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) AddAudio(params *telego.SendAudioParams) {
	u.PlayList.set(params)
}

func (u *User) GetRoomTitle() string {
	return u.Room.Title
}
