package main

import (
	"MusicBot/user"
	"errors"
	"fmt"
	"github.com/mymmrac/telego"
	"os"
)

type Bot struct {
	telegram *telego.Bot
}

func main() {
	botToken := "6210745530:AAGaHIzNOzXlQG9JOMYy1M3DQdxzJ0bjSnY"

	bot := Bot{}

	// Create Bot with debug on
	// Note: Please keep in mind that default logger may expose sensitive information, use in development only
	err := errors.New("")
	bot.telegram, err = telego.NewBot(botToken) //, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get updates channel
	updates, _ := bot.telegram.UpdatesViaLongPolling(nil)

	// Stop reviving updates from update channel
	defer bot.telegram.StopLongPolling()

	admin := user.Init(bot.telegram)

	// Loop through all updates when they came
	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}
		admin.Handler(&update)
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
