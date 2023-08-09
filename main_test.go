package main

import (
	"MusicBot/test/bot"
	testUser "MusicBot/test/user"
	"MusicBot/user"
	"fmt"
	"github.com/mymmrac/telego"
	"os"
	"testing"
)

func testInitTG() (*bot.Bot, <-chan telego.Update) {
	botToken := "6210745530:AAGaHIzNOzXlQG9JOMYy1M3DQdxzJ0bjSnY"

	// Create Bot with debug on
	// Note: Please keep in mind that default logger may expose sensitive information, use in development only
	tg, err := bot.NewBot(botToken) //, telego.WithDefaultDebugLogger())
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

func testStart(updates <-chan telego.Update, admin *user.Admin) {
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

func Test(t *testing.T) {
	tg, updates := testInitTG()

	// Stop reviving updates from update channel
	defer tg.StopLongPolling()

	admin := user.Init(tg)

	tests := []struct {
		id        int64
		firstName string
		username  string
		actions   []testUser.Action
	}{
		{
			id:        101,
			firstName: "receiving",
			username:  "receiving",
			actions: []testUser.Action{
				{
					Send: &telego.Update{
						Message: &telego.Message{Text: "/start"},
					},
				},
				{
					Get: &telego.Update{},
				},
				{
					Send: &telego.Update{},
					Get:  &telego.Update{},
				},
				{
					Send: &telego.Update{},
					Get:  &telego.Update{},
				},
				{
					Send: &telego.Update{},
					Get:  &telego.Update{},
				},
				{
					Send: &telego.Update{},
					Get:  &telego.Update{},
				},
				{
					Send: &telego.Update{},
					Get:  &telego.Update{},
				},
			},
		},
	}

	testStart(updates, admin)
}
