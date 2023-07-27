package main

import (
	"MusicBot/user"
	"fmt"
	"github.com/mymmrac/telego"
	"os"
)

func main() {
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

	// Stop reviving updates from update channel
	defer tg.StopLongPolling()

	admin := user.Init(tg)

	// Loop through all updates when they came
	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}
		admin.Handler(&update)
	}
}
