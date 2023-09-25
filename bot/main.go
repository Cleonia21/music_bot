package main

import (
	"MusicBot/log"
	"MusicBot/telegram"
	"MusicBot/user"
)

func start() {
	// GetMsg updates channel
	updates, err := telegram.TG.UpdatesViaLongPolling(nil)
	if err != nil {
		panic(err)
	}

	// Stop reviving updates from update channel
	defer telegram.TG.StopLongPolling()

	admin := user.Init()

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

	telegram.Init()
	log.Init()
	start()
	//TMain()
}
