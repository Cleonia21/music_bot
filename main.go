package main

import (
	"errors"
	"fmt"
	"github.com/mymmrac/telego"
	"os"
)

type Bot struct {
	telegram *telego.Bot
	audio    *Audio
	users    map[int64]*User
	rooms    map[int64]*Room
}

func (b *Bot) handlingUpdate(update *telego.Update) {
	// Check if update contains a message
	if update.Message == nil {
		return
	}

	audio, err := b.audio.GetParams(update)
	if err != nil {
		b.telegram.Logger().Debugf(err.Error())
		return
	}

	_, err = b.telegram.SendAudio(audio)
	if err != nil {
		b.telegram.Logger().Debugf(err.Error())
		return
	}
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

	bot.audio = CreateAudio()

	// Loop through all updates when they came
	for update := range updates {
		//update.Message.From.ID
		bot.handlingUpdate(&update)
	}
}
