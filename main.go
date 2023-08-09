package main

import (
	"MusicBot/user"
	"fmt"
	"github.com/mymmrac/telego"
	"os"
)

func initTG() (*telego.Bot, <-chan telego.Update) {
	botToken := "6210745530:AAGaHIzNOzXlQG9JOMYy1M3DQdxzJ0bjSnY"

	// Create Bot with debug on
	// Note: Please keep in mind that default logger may expose sensitive information, use in development only
	tg, err := telego.NewBot(botToken) //, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get updates channel
	updates, err := tg.UpdatesViaLongPolling(nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return tg, updates
}

func start() {
	tg, updates := initTG()

	// Stop reviving updates from update channel
	defer tg.StopLongPolling()

	admin := user.Init(tg)

	//-809440484
	// Loop through all updates when they came
	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}
		if update.Message != nil && update.Message.Chat.Type != "private" {
			continue
		}
		admin.Handler(&update)
	}
}

func main() {
	start()
}
