package usecase

import (
	"music_bot/internal/entity"
	"strconv"
	"strings"
)

type UnregUser struct {
	user   *entity.UnregUser
	pass   string
	hostId entity.UserID
	sender Sender
}

func (u *UnregUser) setUpdate(update entity.Update) (action string) {
	switch update.Command {
	case "/start":
		u.sender.SendStartMenu(u.user.ID)
	case "/menu":
		u.sender.SendStartMenu(u.user.ID)
	case "/info":
		u.sender.SendStartMenu(u.user.ID)
	case "/host_role":
		return "host"
	case "/child_role":
		u.sender.ChildRegGreeting(u.user.ID)
	default:
		return u.parseSecretMsg(update.Text)
	}
	return ""
}

func (u *UnregUser) setAction(action string) {
	switch action {
	case "host not found":
		u.sender.UnidentSecretMsg(u.user.ID)
	case "pass incorrect":
		u.sender.UnidentSecretMsg(u.user.ID)
	default:
		u.sender.UnknownError(u.user.ID)
	}
}

func (u *UnregUser) parseSecretMsg(text string) (action string) {
	strs := strings.Split(text, "/")
	if len(strs) != 3 {
		u.sender.UnidentSecretMsg(u.user.ID)
		return
	}
	if strs[0] != "secretMessage" {
		u.sender.UnidentSecretMsg(u.user.ID)
		return
	}
	hostId, err := strconv.Atoi(strs[1][1:])
	if err != nil {
		u.sender.UnidentSecretMsg(u.user.ID)
		return
	}
	u.hostId = entity.NewUserId(hostId)

	u.pass = strs[2]
	return "child"
}
