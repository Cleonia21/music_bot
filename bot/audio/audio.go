package audio

import (
	"MusicBot/audio/fromURL"
	"MusicBot/audio/yandex"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"sync"
)

type Audio struct {
	yandex yandex.Music
	mutex  sync.Mutex
}

func Init() *Audio {
	a := Audio{
		yandex: yandex.Music{},
	}
	a.yandex.Authorization()
	return &a
}

func (a *Audio) GetParams(update *telego.Update) (*telego.SendAudioParams, error) {
	var audio *telego.SendAudioParams

	if update.Message.Audio != nil {
		audio = tu.Audio(tu.ID(update.Message.Chat.ID), telego.InputFile{FileID: update.Message.Audio.FileID})
		audio.WithTitle(update.Message.Audio.Title).WithPerformer(update.Message.Audio.Performer)
	} else {
		a.mutex.Lock()
		params, err := a.yandex.URLtoAudioParams(update.Message.Text)
		if err != nil {
			return nil, err
		}

		audio, err = fromURL.FromURL(params)
		a.mutex.Unlock()
	}
	return audio, nil
}
