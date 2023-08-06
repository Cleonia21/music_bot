package utils

import "github.com/mymmrac/telego"

func UpdateToID(update *telego.Update) (id telego.ChatID) {
	var user *telego.User
	if update.Message != nil {
		user = update.Message.From
	} else if update.CallbackQuery != nil {
		user = &update.CallbackQuery.From
	}
	id.ID = user.ID
	id.Username = user.Username
	return
}
