package yandex

import (
	"context"
	"fmt"
	"github.com/ndrewnee/go-yamusic/yamusic"
)

type PlayList struct {
	playList yamusic.PlaylistsResult
}

func (m *Music) CreatePlayList(title string) (playList yamusic.PlaylistsResult) {
	createResp, _, err := m.client.Playlists().Create(
		context.Background(),
		title,
		true,
	)

	res := createResp.Result

	fmt.Printf("create result: %v err: %v", res, err)

	//https://music.yandex.com/album/23345073/track/106849229
	addResp, _, err := m.client.Playlists().AddTracks(
		context.Background(),
		res.Kind,
		1,
		[]yamusic.PlaylistsTrack{
			{
				ID:      106849229,
				AlbumID: 23345073,
			},
		},
		nil,
	)

	res = addResp.Result

	fmt.Printf("add result: %v err: %v", res, err)
	return yamusic.PlaylistsResult{}
}
