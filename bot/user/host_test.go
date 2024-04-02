package user

import (
	"MusicBot/audio"
	"MusicBot/user/utils"
	"github.com/mymmrac/telego"
	"testing"
)

func Test_hostUser_init(t *testing.T) {
	var user hostUser

	chatID := utils.NewUserID(10, "10", "10")
	a := &audio.Audio{}
	ch := make(chan telego.Update)

	user.init(chatID, a, ch)

	emptyUserID := utils.UserID{}
	if user.id == emptyUserID ||
		user.getFromTgCh == nil ||
		user.pass == "" ||
		user.audio == nil ||
		user.getFromUsersCh == nil ||
		user.sendToUsersChs == nil {
		t.Error("the field of structure hostUser has an incorrect value")
	}
}

func Test_hostUser_join(t *testing.T) {
	var HostUser hostUser

	HostChatID := utils.NewUserID(10, "10", "10")
	HostAudio := &audio.Audio{}
	TgCh := make(chan telego.Update)

	HostUser.init(HostChatID, HostAudio, TgCh)

	JoinChatID := utils.NewUserID(20, "20", "20")
	JoinCH := make(chan msgBetweenUsers)

	hostInCh := HostUser.join(JoinChatID, JoinCH)
	if hostInCh != HostUser.getFromUsersCh {
		t.Errorf("join() = %v, want %v", hostInCh, HostUser.getFromUsersCh)
	}
	if HostUser.sendToUsersChs[JoinChatID] != JoinCH {
		t.Errorf("join(): the joining user is not recorded correctly")
	}
}

func Test_hostUser_out(t *testing.T) {
	var user hostUser
	chatID := utils.NewUserID(10, "10", "10")
	a := &audio.Audio{}
	ch := make(chan telego.Update)
	user.init(chatID, a, ch)

	joiningUsersParam := []struct {
		id         utils.UserID
		senderInCh chan msgBetweenUsers
	}{
		{
			utils.NewUserID(1, "1", "1"),
			make(chan msgBetweenUsers),
		},
		{
			utils.NewUserID(2, "2", "2"),
			make(chan msgBetweenUsers),
		},
		{
			utils.NewUserID(3, "3", "3"),
			make(chan msgBetweenUsers),
		},
	}
	for _, param := range joiningUsersParam {
		user.join(param.id, param.senderInCh)
	}

	user.out()

	for _, param := range joiningUsersParam {
		msg := <-param.senderInCh
		if msg.id != "out" && msg.from != user.id {
			t.Errorf("out func work incorrect")
		}
	}
}
