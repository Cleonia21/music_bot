package playList

import (
	"github.com/mymmrac/telego"
	"testing"
)

func TestPlayList(t *testing.T) {
	ids := []telego.ChatID{
		{Username: "1"},
		{Username: "2"},
		{Username: "3"},
		{Username: "4"},
	}

	audioParams := []*telego.SendAudioParams{
		{Title: "0"},
		{Title: "1"},
		{Title: "2"},
		{Title: "3"},
		{Title: "4"},
		{Title: "5"},
		{Title: "6"},
		{Title: "7"},
		{Title: "8"},
		{Title: "9"},
	}

	type setData struct {
		id    telego.ChatID
		audio *telego.SendAudioParams
	}
	type getData struct {
		audios []*telego.SendAudioParams
		errs   []string
	}
	tests := []struct {
		name    string
		setData map[int]setData
		getData []getData
		summary []Summary
	}{
		{
			name: "",
			setData: map[int]setData{
				0: {
					id:    telego.ChatID{Username: "1"},
					audio: &telego.SendAudioParams{Title: "11"},
				},
				1: {
					id:    telego.ChatID{Username: "1"},
					audio: &telego.SendAudioParams{Title: "12"},
				},
				2: {
					id:    telego.ChatID{Username: "1"},
					audio: &telego.SendAudioParams{Title: "13"},
				},
				3: {
					id:    telego.ChatID{Username: "2"},
					audio: &telego.SendAudioParams{Title: "21"},
				},
				4: {
					id:    telego.ChatID{Username: "3"},
					audio: &telego.SendAudioParams{Title: "31"},
				},
				5: {
					id:    telego.ChatID{Username: "3"},
					audio: &telego.SendAudioParams{Title: "32"},
				},
			},
			getData: []getData{
				{
					audios: nil,
					errs:   nil,
				},
			},
			summary: []Summary{
				{
					ID:  telego.ChatID{},
					Num: 0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PlayList{}
			p.Init()

			if err := p.SetAudio(tt.args.id, tt.args.audio); (err != nil) != tt.wantErr {
				t.Errorf("SetAudio() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
