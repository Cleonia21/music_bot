package playList

import (
	"github.com/mymmrac/telego"
	"reflect"
	"sort"
	"testing"
)

func TestPlayList(t *testing.T) {
	type wantData struct {
		audios []*telego.SendAudioParams
		errs   []string
	}

	p := new(PlayList)
	p.Init()

	_ = p.SetAudio(telego.ChatID{Username: "1"}, &telego.SendAudioParams{Title: "11"})
	_ = p.SetAudio(telego.ChatID{Username: "1"}, &telego.SendAudioParams{Title: "12"})
	_ = p.SetAudio(telego.ChatID{Username: "1"}, &telego.SendAudioParams{Title: "13"})
	_ = p.SetAudio(telego.ChatID{Username: "2"}, &telego.SendAudioParams{Title: "21"})
	_ = p.SetAudio(telego.ChatID{Username: "3"}, &telego.SendAudioParams{Title: "31"})
	_ = p.SetAudio(telego.ChatID{Username: "3"}, &telego.SendAudioParams{Title: "32"})

	tests := []struct {
		name        string
		wantData    wantData
		wantSummary []Summary
	}{
		{
			name: "",
			wantData: wantData{
				audios: []*telego.SendAudioParams{
					{Title: "11"},
					{Title: "21"},
					{Title: "31"},
				},
				errs: nil,
			},
			wantSummary: []Summary{
				{
					ID:  telego.ChatID{Username: "1"},
					Num: 3,
				},
				{
					ID:  telego.ChatID{Username: "2"},
					Num: 1,
				},
				{
					ID:  telego.ChatID{Username: "3"},
					Num: 2,
				},
			},
		},
		{
			name: "",
			wantData: wantData{
				audios: []*telego.SendAudioParams{
					{Title: "12"},
					{Title: "32"},
				},
				errs: []string{
					"2",
				},
			},
			wantSummary: []Summary{
				{
					ID:  telego.ChatID{Username: "1"},
					Num: 2,
				},
				{
					ID:  telego.ChatID{Username: "2"},
					Num: 0,
				},
				{
					ID:  telego.ChatID{Username: "3"},
					Num: 1,
				},
			},
		},
		{
			name: "",
			wantData: wantData{
				audios: []*telego.SendAudioParams{
					{Title: "13"},
				},
				errs: []string{
					"2",
					"3",
				},
			},
			wantSummary: []Summary{
				{
					ID:  telego.ChatID{Username: "1"},
					Num: 1,
				},
				{
					ID:  telego.ChatID{Username: "2"},
					Num: 0,
				},
				{
					ID:  telego.ChatID{Username: "3"},
					Num: 0,
				},
			},
		},
		{
			name: "",
			wantData: wantData{
				audios: nil,
				errs: []string{
					"1",
					"2",
					"3",
				},
			},
			wantSummary: []Summary{
				{
					ID:  telego.ChatID{Username: "1"},
					Num: 0,
				},
				{
					ID:  telego.ChatID{Username: "2"},
					Num: 0,
				},
				{
					ID:  telego.ChatID{Username: "3"},
					Num: 0,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getSummary := p.GetSummary()

			sort.Slice(getSummary, func(i, j int) bool {
				return getSummary[i].ID.Username < getSummary[j].ID.Username
			})

			if !reflect.DeepEqual(getSummary, tt.wantSummary) {
				t.Errorf("get summary = %v, want summary = %v", getSummary, tt.wantSummary)
			}

			getAudios, getErrors := p.GetAudio()
			sort.Slice(getAudios, func(i, j int) bool {
				return getAudios[i].Title < getAudios[j].Title
			})

			sort.Slice(getErrors, func(i, j int) bool {
				return getErrors[i] < getErrors[j]
			})

			if !reflect.DeepEqual(getAudios, tt.wantData.audios) {
				t.Errorf("get audios = %v, want audios = %v", getAudios, tt.wantData.audios)
			}

			if !reflect.DeepEqual(getErrors, tt.wantData.errs) {
				t.Errorf("get errors = %v, want errors = %v", getErrors, tt.wantData.errs)
			}

		})
	}
}
