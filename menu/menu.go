package menu

import (
	"MusicBot/users"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

/*
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
*/

const (
	StateStart  = 0
	StateNoRoom = 1

	StateCreateRoomWaitTitle = 20
	StateCreateRoomWaitPass  = 21

	StateJoinRoomTitle = 30
	StateJoinRoomPass  = 31

	StateRoomRoot  = 4
	StateRoomGuest = 5

	StateIncorrectAnswer = 100
)

type Menu struct {
	keyboards map[int]*telego.ReplyKeyboardMarkup
	texts     map[int]string

	prefix string
}

func (m *Menu) Init() {
	m.keyboards = map[int]*telego.ReplyKeyboardMarkup{
		StateNoRoom: tu.Keyboard(
			tu.KeyboardRow( // Row 1
				tu.KeyboardButton("1"), // Column 1
				tu.KeyboardButton("2"), // Column 2
			),
		).WithResizeKeyboard(), //.WithInputFieldPlaceholder("Select something")
	}
	m.texts = map[int]string{
		StateStart: "",
		StateNoRoom: "1. Создать комнату\n" +
			"2. Присоедениться к комнате",
		StateCreateRoomWaitTitle: "Введите название",
		StateCreateRoomWaitPass:  "Введите пароль",
	}
}

func (m *Menu) getKeyboard(state int) telego.ReplyMarkup {
	keyboard := m.keyboards[state]
	if keyboard == nil {
		return tu.ReplyKeyboardRemove()
	}
	return keyboard
}

func (m *Menu) Get(user *users.User) *telego.SendMessageParams {
	var text string
	if m.prefix != "" {
		text = m.prefix
	} else {
		text = m.texts[user.State]
	}
	message := tu.Message(user.ChatID, text)
	message.WithReplyMarkup(m.getKeyboard(user.State))

	m.prefix = ""

	return message
}

const (
	PrefixNoSuchAnswer = 1
)

func (m *Menu) SetPrefix(prefix int) {
	switch prefix {
	case PrefixNoSuchAnswer:
		m.prefix = "Нет такого варианта ответа"
	}
}
