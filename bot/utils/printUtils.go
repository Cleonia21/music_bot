package utils

import (
	"fmt"
	"github.com/mymmrac/telego"
)

type Utils struct {
}

func UpdateToStr(update *telego.Update) string {
	if update == nil {
		return ""
	}
	if update.Message != nil {
		return MsgToStr(update.Message)
	} else if update.CallbackQuery != nil {
		query := update.CallbackQuery
		return fmt.Sprintf("{CQ: username(%v)text(%v)data(%v)}",
			query.From.Username,
			query.Message.Text,
			query.Data,
		)
	} else {
		return ""
	}
}

func MsgToStr(msg *telego.Message) string {
	if msg != nil {
		return fmt.Sprintf("{MSG: username(%v)text(%v)btns(%v)audio(%v)}",
			msg.Chat.Username,
			textCompression(msg.Text, 20),
			replyMarkupToStr(msg.ReplyMarkup),
			audioToStr(msg.Audio),
		)
	}
	return "{MSG: nil}"
}

func textCompression(text string, length int) string {
	if len(text) > length {
		return text[:length] + "..."
	}
	return text
}

func replyMarkupToStr(markup *telego.InlineKeyboardMarkup) string {
	if markup == nil {
		return "nil"
	}
	btns := "{BTNS: "
	for _, row := range markup.InlineKeyboard {
		for _, btn := range row {
			btns += fmt.Sprintf("{text(%v),data(%v)}", btn.Text, btn.CallbackData)
		}
	}
	return btns + "}"
}

func audioToStr(audio *telego.Audio) string {
	if audio == nil {
		return "nil"
	}
	return fmt.Sprintf("{AUDIO: title(%v)perf(%v)}", audio.Title, audio.Performer)
}
