package main

import (
	"MusicBot/test/bot"
	testUser "MusicBot/test/user"
	"MusicBot/user"
	"fmt"
	"github.com/mymmrac/telego"
	"os"
	"sync"
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

	// GetMsg updates channel
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

func TMain() {
	tg, updates := testInitTG()

	// Stop reviving updates from update channel
	defer tg.StopLongPolling()

	admin := user.Init()

	tests := []struct {
		id      int64
		actions []testUser.Action
	}{
		{
			id: 101,
			actions: []testUser.Action{
				{
					Send: "/start",
				},
				{
					AnswerToMsg: "host",
				},
				{
					UnprocGet: 4,
				},
				{
					Send: "https://music.yandex.ru/album/19435876/track/95386895",
				},
				{
					GetText: "Отправил в очередь",
				},
				{
					Send: "https://music.yandex.ru/album/19435876/track/95386889",
				},
				{
					GetText: "Отправил в очередь",
				},
				{
					Send: "https://music.yandex.ru/album/19435876/track/95386880",
				},
				{
					GetText: "Отправил в очередь",
				},
				{
					Send: "https://music.yandex.ru/album/19435876/track/95386889",
				},
				{
					GetText: "Отправил в очередь",
				},
			},
		},
		{
			id: 201,
			actions: []testUser.Action{
				{
					Send: "/start",
				},
				{
					AnswerToMsg: "send",
				},
				{
					UnprocGet: 1,
				},
				{
					Send: "https://music.yandex.ru/album/3280462/track/27395628",
				},
				{
					GetText: "Отправил в очередь",
				},
				{
					Send: "https://music.yandex.ru/album/12599890/track/72920590",
				},
				{
					GetText: "Отправил в очередь",
				},
				{
					Send: "https://music.yandex.ru/album/18700169/track/93015961",
				},
				{
					GetText: "Отправил в очередь",
				},
				{
					Send: "https://music.yandex.ru/album/4937815/track/35590676",
				},
				{
					GetText: "Отправил в очередь",
				},
			},
		},
	}

	server := tg.Server
	var wg sync.WaitGroup
	wg.Add(len(tests))

	for _, userTestCase := range tests {
		u := server.AddUser(userTestCase.id)
		go u.SetScript(userTestCase.actions, &wg)
	}

	go server.Start()
	go testStart(updates, admin)
	wg.Wait()
}
