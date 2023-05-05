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

type Bot struct {
	telegram  *telego.Bot
	audio     *Audio
	users     map[int64]*users.User
	roomsRoot map[string]*users.User
	menu      *menu.Menu
}

func (b *Bot) newUser(update *telego.Update) *users.User {
	user := b.users[update.Message.From.ID]
	if user == nil {
		user = &users.User{}
		user.Constructor(
			tu.ID(update.Message.Chat.ID),
			update.Message.From.ID,
			update.Message.From.Username,
		)
	}
	b.users[user.ID] = user
	return user
}

func (b *Bot) stateJoinRoomWaitPass(user *users.User, update *telego.Update) error {
	title := user.MessHistory[0]
	pass := update.Message.Text

	roomRoot := b.roomsRoot[title]
	if roomRoot == nil {
		return errors.New("title not found")
	}
	err := roomRoot.AddInRoom(user, pass)
	if err == nil {
		user.State = menu.StateRoomGuest
	}
	return err
}

func (b *Bot) stateJoinRoomWaitTitle(user *users.User, update *telego.Update) {
	user.MessHistory[0] = update.Message.Text
	user.State = menu.StateJoinRoomWaitPass
}

func (b *Bot) stateCreateRoomWaitPass(user *users.User, update *telego.Update) {
	title := user.MessHistory[0]
	pass := update.Message.Text
	user.CreateRoom(title, pass)
	b.roomsRoot[title] = user
	user.State = menu.StateRoomRoot
}

func (b *Bot) stateCreateRoomWaitTitle(user *users.User, update *telego.Update) {
	user.MessHistory[0] = update.Message.Text
	user.State = menu.StateCreateRoomWaitPass
}

func (b *Bot) stateNoRoom(user *users.User, update *telego.Update) {
	switch update.Message.Text {
	case "1":
		user.State = menu.StateCreateRoomWaitTitle
	case "2":
		user.State = menu.StateJoinRoomWaitTitle
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
	case menu.StateJoinRoomWaitTitle:
		b.stateJoinRoomWaitTitle(user, update)
	case menu.StateJoinRoomWaitPass:
		b.stateJoinRoomWaitPass()

	}
	_, _ = b.telegram.SendMessage(b.menu.Get(user))
}

func main() {
	botToken := TOKEN

	bot := Bot{}
	bot.users = make(map[int64]*users.User)
	bot.roomsRoot = make(map[string]*users.User)

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
