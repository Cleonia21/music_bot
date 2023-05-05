package main

import (
	"MusicBot/menu"
	"MusicBot/users"
	"errors"
	"fmt"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"os"
)

//const (
//	stateStart  = 0
//	stateNoRoom = 1
//
//	stateCreateRoom = 2
//	stateJoinRoom   = 3
//
//	stateRoomRoot  = 4
//	stateRoomGuest = 5
//)

type Bot struct {
	telegram *telego.Bot
	audio    *Audio
	users    map[int64]*users.User
	rooms    map[string]*users.Room
	menu     *menu.Menu
}

func (b *Bot) joinRoom(user *users.User, title, pass string) error {
	if user.Room != nil {
		// correctly out from room
		user.Room.Dell(b.users[user.ID])
	}
	Room, ok := b.rooms[title]
	if ok {
		return Room.Join(b.users[user.ID], pass)
	}
	return errors.New("incorrect title")
}

func (b *Bot) createRoom(user *users.User, title, pass string) {
	var room users.Room
	room.CreateRoom(
		user,
		title,
		pass,
	)
	b.rooms[title] = &room
}

func (b *Bot) newUser(update *telego.Update) *users.User {
	userID := update.Message.From.ID
	user, ok := b.users[userID]
	if ok {
		return user
	}
	// check for duplicates
	user = users.Constructor(
		tu.ID(update.Message.Chat.ID),
		update.Message.From.ID,
		update.Message.From.Username,
	)
	b.users[user.ID] = user
	return user
}

/*
func (b *Bot) stateJoinRoom(user *users.User, update *telego.Update) string {
	title, ok := user.MessHistory[0]
	if !ok {
		user.MessHistory[0] = update.Message.Text
		return menu.JoinRoomPass()
	}
	// clear user message history
	pass := update.Message.Text
	err := b.joinRoom(user, title, pass)
	if err != nil {
		user.State = stateNoRoom
		return menu.RoomNotCreated() + "\n" + menu.RoomAction()
	}
	user.State = stateRoomGuest
	return menu.JoinRoomComp()
}
*/

func (b *Bot) stateCreateRoomWaitTitle(user *users.User, update *telego.Update) {
	user.MessHistory[0] = update.Message.Text
	user.State = menu.StateCreateRoomWaitPass
}

func (b *Bot) stateCreateRoomWaitPass(user *users.User, update *telego.Update) {
	title := user.MessHistory[0]
	pass := update.Message.Text
	b.createRoom(user, title, pass)
	user.State = menu.StateRoomRoot
}

func (b *Bot) stateNoRoom(user *users.User, update *telego.Update) {
	switch update.Message.Text {
	case "1":
		user.State = menu.StateCreateRoomWaitTitle
	case "2":
		user.State = menu.StateJoinRoomTitle
	default:
		b.menu.SetPrefix(menu.PrefixNoSuchAnswer)
	}
}

func (b *Bot) stateStart(user *users.User) {
	user.State = menu.StateNoRoom
}

func (b *Bot) handlingUpdate(update *telego.Update) {
	if update.Message == nil {
		return
	}

	user := b.newUser(update)
	switch user.State {
	case menu.StateStart:
		b.stateStart(user)
	case menu.StateNoRoom:
		b.stateNoRoom(user, update)
	case menu.StateCreateRoomWaitTitle:
		b.stateCreateRoomWaitTitle(user, update)
	case menu.StateCreateRoomWaitPass:
		b.stateCreateRoomWaitPass(user, update)

		//case menu.StateJoinRoom:
		//	b.stateJoinRoom(user, update)
	}
	_, _ = b.telegram.SendMessage(b.menu.Get(user))
}

func main() {
	botToken := TOKEN

	bot := Bot{}
	bot.users = make(map[int64]*users.User)
	bot.rooms = make(map[string]*users.Room)

	bot.menu = &menu.Menu{}
	bot.menu.Init()

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
