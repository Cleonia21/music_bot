package playList

import (
	"github.com/mymmrac/telego"
)

type PlayList struct {
	audios map[telego.ChatID]*queue
}

func (p *PlayList) Init() {
	p.audios = make(map[telego.ChatID]*queue)
}

func (p *PlayList) SetAudio(id telego.ChatID, audio *telego.SendAudioParams) error {
	qChan, ok := p.audios[id]
	if !ok {
		qChan = initQueue()
		p.audios[id] = qChan
	}
	err := qChan.set(audio)
	return err
}

func (p *PlayList) GetAudio() (audios []*telego.SendAudioParams, errs []string) {
	for id, qe := range p.audios {
		audio, err := qe.get()
		if err != nil {
			errs = append(errs, id.Username)
		} else {
			audios = append(audios, audio)
		}
	}
	return
}

type Summary struct {
	ID  telego.ChatID
	Num int
}

func (p *PlayList) GetSummary() (s []Summary) {
	for id, qe := range p.audios {
		s = append(s, Summary{ID: id, Num: qe.len})
	}
	return
}

func (p *PlayList) UserTrackNum(id telego.ChatID) int {
	que, ok := p.audios[id]
	if !ok {
		return 0
	}
	return que.len
}
