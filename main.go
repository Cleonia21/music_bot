package main

import (
	"errors"
	"fmt"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"os"
)

const (
	stateStart  = 0
	stateNoRoom = 1

	stateCreateRoom = 2
	stateJoinRoom   = 3

	stateRoomRoot  = 4
	stateRoomGuest = 5
)

type Bot struct {
	telegram *telego.Bot
	audio    *Audio
	users    map[int64]*User
	rooms    map[string]*Room
}

func (b *Bot) joinRoom(user *User, title, pass string) error {
	if user.Room != nil {
		// correctly out from room
		user.Room.Dell(b.users[user.ID])
	}
	room, ok := b.rooms[title]
	if ok {
		return room.Join(b.users[user.ID], pass)
	}
	return errors.New("incorrect title")
}

func (b *Bot) createRoom(user *User, title, pass string) {
	var room Room
	room.CreateRoom(
		user,
		title,
		pass,
	)
	b.rooms[title] = &room
}

func (b *Bot) newUser(update *telego.Update) *User {
	user := CreateUser(
		tu.ID(update.Message.Chat.ID),
		update.Message.From.ID,
		update.Message.From.Username,
	)
	b.users[user.ID] = user
	return user
}

func (b *Bot) stateJoinRoom(user *User, update *telego.Update) string {
	title, ok := user.MessHistory[0]
	if !ok {
		user.MessHistory[0] = update.Message.Text
		return menu{}.createRoomPass()
	}
	// clear user message history
	pass := update.Message.Text
	err := b.joinRoom(user, title, pass)
	if err != nil {
		return menu{}.joinRoomTitle()
	}
	user.State = stateRoomGuest
	return menu{}.joinRoomComp()
}

func (b *Bot) stateCreateRoom(user *User, update *telego.Update) string {
	title, ok := user.MessHistory[0]
	if !ok {
		user.MessHistory[0] = update.Message.Text
		return menu{}.createRoomPass()
	}
	// clear user message history
	pass := update.Message.Text
	b.createRoom(user, title, pass)
	user.State = stateRoomRoot
	return menu{}.roomCreated()
}

func (b *Bot) stateNoRoom(user *User, update *telego.Update) string {
	switch update.Message.Text {
	case "1":
		user.State = stateCreateRoom
		return menu{}.createRoomTitle()
	case "2":
		user.State = stateJoinRoom
		return menu{}.joinRoomTitle()
	default:
		return menu{}.noSuchAnswer()
	}
}

func (b *Bot) stateStart(user *User) string {
	user.State = stateNoRoom
	return menu{}.roomAction()
}

func (b *Bot) handlingUpdate(update *telego.Update) {
	if update.Message == nil {
		return
	}

	user := b.newUser(update)
	var text string
	switch user.State {
	case stateStart:
		text = b.stateStart(user)
	case stateNoRoom:
		text = b.stateNoRoom(user, update)
	case stateCreateRoom:
		text = b.stateCreateRoom(user, update)
	}
	_, _ = b.telegram.SendMessage(tu.Message(user.ChatID, text))
}

func main() {
	botToken := TOKEN

	bot := Bot{}
	// Create Bot with debug on
	// Note: Please keep in mind that default logger may expose sensitive information, use in development only
	err := errors.New("")
	bot.telegram, err = telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get updates channel
	updates, _ := bot.telegram.UpdatesViaLongPolling(nil)

	// Stop reviving updates from update channel
	defer bot.telegram.StopLongPolling()

	bot.audio = CreateAudio(bot.telegram)

	// Loop through all updates when they came
	for update := range updates {
		//update.Message.From.ID
		bot.handlingUpdate(&update)
	}
}

/*
func main() {
	botToken := TOKEN

	// Note: Please keep in mind that default logger may expose sensitive information, use in development only
	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Keyboard parameters
	keyboard := tu.Keyboard(
		tu.KeyboardRow( // Row 1
			tu.KeyboardButton("Button 1"), // Column 1
			tu.KeyboardButton("Button 2"), // Column 2
		),
		tu.KeyboardRow( // Row 2
			tu.KeyboardButton("Contact").WithRequestContact(),   // Column 1
			tu.KeyboardButton("Location").WithRequestLocation(), // Column 2
		),
		tu.KeyboardRow( // Row 3
			tu.KeyboardButton("Poll Any").WithRequestPoll(tu.PollTypeAny()),         // Column 1
			tu.KeyboardButton("Poll Regular").WithRequestPoll(tu.PollTypeRegular()), // Column 2
			tu.KeyboardButton("Poll Quiz").WithRequestPoll(tu.PollTypeQuiz()),       // Column 3
		),
	).WithResizeKeyboard().WithInputFieldPlaceholder("Select something")

	// Get updates channel
	updates, _ := bot.UpdatesViaLongPolling(nil)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		// Message parameters
		message := tu.Message(
			tu.ID(update.Message.Chat.ID),
			"My message",
		).WithReplyMarkup(keyboard)
		// Sending message
		_, _ = bot.SendMessage(message)
	}

}

*/
