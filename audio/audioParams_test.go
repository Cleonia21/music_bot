package audio

import (
	"github.com/mymmrac/telego"
	"reflect"
	"testing"
)

func TestAudio_GetParams(t *testing.T) {
	type args struct {
		update *telego.Update
	}
	tests := []struct {
		name    string
		args    args
		want    *telego.SendAudioParams
		wantErr bool
	}{
		{
			"",
			args{
				&telego.Update{
					Message: &telego.Message{
						Chat: telego.Chat{ID: 0},
						Text: "https://music.yandex.ru/album/8958861/track/55049473",
					},
				},
			},
			&telego.SendAudioParams{},
			false,
		},
	}
	a := Init()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := a.GetParams(tt.args.update)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetParams() got = %v, want %v", got, tt.want)
			}
		})
	}
}
