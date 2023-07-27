package playList

import (
	"errors"
	"github.com/mymmrac/telego"
)

const queueMaxLen = 20

type queue struct {
	len   int
	queue chan *telego.SendAudioParams
}

func initQueue() *queue {
	q := queue{
		len:   0,
		queue: make(chan *telego.SendAudioParams, 20),
	}
	return &q
}

func (q *queue) set(audio *telego.SendAudioParams) error {
	if q.len == queueMaxLen {
		return errors.New("в очередь добавлено максимально количество треков")
	}
	q.len++
	q.queue <- audio
	return nil
}

func (q *queue) get() (audio *telego.SendAudioParams, err error) {
	if q.len == 0 {
		return nil, errors.New("пользователь не добавлял треков")
	}
	audio = <-q.queue
	q.len--
	return
}
