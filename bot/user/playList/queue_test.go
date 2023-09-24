package playList

import (
	"github.com/mymmrac/telego"
	"testing"
)

func Test_queue(t *testing.T) {
	tests := []struct {
		name   string
		audios []*telego.SendAudioParams
	}{
		{
			audios: []*telego.SendAudioParams{
				{Title: "1"},
			},
		},
		{
			audios: []*telego.SendAudioParams{
				{Title: "1"},
				{Title: "2"},
				{Title: "3"},
				{Title: "4"},
			},
		},
		{
			audios: []*telego.SendAudioParams{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := initQueue()

			for i, audio := range tt.audios {
				err := q.set(audio)
				if err != nil {
					t.Errorf(err.Error())
				}
				if q.len != i+1 {
					t.Errorf("want que len = %v, get que len = %v", i+1, q.len)
				}
			}

			for _, wantAudio := range tt.audios {
				getAudio, err := q.get()
				if err != nil {
					t.Errorf(err.Error())
				}
				if getAudio.Title != wantAudio.Title {
					t.Errorf("get = %v, want = %v", getAudio, wantAudio)
				}
			}
		})
	}
}
