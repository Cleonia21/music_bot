package telegram

import (
	"github.com/mymmrac/telego"
)

var TG *telego.Bot

func Init() {
	botToken := "6210745530:AAGaHIzNOzXlQG9JOMYy1M3DQdxzJ0bjSnY"

	// Create Bot with debug on
	// Note: Please keep in mind that default logger may expose sensitive information, use in development only
	var err error
	TG, err = telego.NewBot(botToken) //, telego.WithDefaultDebugLogger())
	if err != nil {
		panic(err)
	}
}
