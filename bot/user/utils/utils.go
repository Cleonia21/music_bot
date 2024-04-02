package utils

import (
	"fmt"
	"github.com/mymmrac/telego"
)

type UserID struct {
	ChatID    telego.ChatID
	ID        int64
	Username  string
	FirstName string
}

func NewUserID(ID int64, Username string, FirstName string) UserID {
	id := UserID{ChatID: telego.ChatID{ID, Username}, ID: ID, Username: Username, FirstName: FirstName}
	return id
}

func UpdateToID(update *telego.Update) (id UserID) {
	var user *telego.User
	if update.Message != nil {
		user = update.Message.From
	} else if update.CallbackQuery != nil {
		user = &update.CallbackQuery.From
	}
	id.ChatID.ID = user.ID
	id.ID = user.ID
	id.ChatID.Username = user.Username
	id.Username = user.Username
	id.FirstName = user.FirstName
	return
}

func UserNameInserting(before string, id UserID, after string) string {
	return fmt.Sprintf("%v<a href=\"tg://user?id=%v\">%v</a>%v", before, id.ID, id.FirstName, after)
}
