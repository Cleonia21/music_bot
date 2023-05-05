package users

import "github.com/mymmrac/telego"

type playList struct {
	audios []*telego.SendAudioParams
}

func (p *playList) constructor() {
	p.audios = make([]*telego.SendAudioParams, 0, 10)
}

func (p *playList) set(audio *telego.SendAudioParams) {
	p.audios = append(p.audios, audio)
}

func (p *playList) get() *telego.SendAudioParams {
	audio := p.audios[0]
	p.audios[0] = nil
	p.audios = p.audios[1:]
	return audio
}
