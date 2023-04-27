package music

import (
	"context"
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/ndrewnee/go-yamusic/yamusic"
	"github.com/rubyist/circuitbreaker"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Music struct {
	client *yamusic.Client
}

func (m *Music) Authorization() {
	// constructing http client with circuit breaker
	// it implements yamusic.Doer interface
	circuitClient := circuit.NewHTTPClient(time.Second*5, 10, nil)
	m.client = yamusic.NewClient(
		// if you want http client with circuit breaker
		yamusic.HTTPClient(circuitClient),
		// provide user_id and access_token (needed by some methods)
		yamusic.AccessToken(352880607, "y0_AgAAAAAVCIffAAG8XgAAAADZ9NQjGLw1kvrFRHiWHhRVb3UKD1ShpsA"),
	)
}

func (m *Music) CreatedPlaylist(title string) {
	// create new public playlist. Need access token
	_, _, err := m.client.Playlists().Create(context.Background(), title, true)
	if err != nil {
		log.Fatal(err)
	}
	//log.Println("Created playlist: ", createdPlaylist)
}

func (m *Music) GetPlayList(title string) (yamusic.PlaylistsResult, error) {
	lists, _, err := m.client.Playlists().List(context.Background(), m.client.UserID())
	if err != nil {
		return yamusic.PlaylistsResult{}, err
	}
	for _, list := range lists.Result {
		if list.Title == title {
			return list, nil
		}
	}
	return yamusic.PlaylistsResult{}, errors.New("playlist not found")
}

func (m *Music) SearchTrack(id int) {
	tracks, _, err := m.client.Tracks().Get(context.Background(), id)
	if err != nil {
		return
	}
	for _, name := range tracks.Result {
		fmt.Println(name.Title)
	}
}

// https://music.yandex.ru/album/17678543/track/89854895
func (m *Music) parseTrackURL(trackURL string) (id int, albumID int, err error) {
	matched, err := regexp.MatchString(`^https:\/\/music\.yandex\.ru\/album\/[0-9]*\/track\/[0-9]*`, trackURL)
	if err != nil {
		return 0, 0, err
	}
	if matched != true {
		return id, albumID, errors.New("invalid track url")
	}
	parse, err := url.Parse(trackURL)
	if err != nil {
		return 0, 0, err
	}
	p := strings.Split(parse.Path, "/")
	id, err = strconv.Atoi(p[4])
	if err != nil {
		return 0, 0, err
	}
	albumID, err = strconv.Atoi(p[2])
	if err != nil {
		return 0, 0, err
	}
	//fmt.Println(id, albumID)
	if id == 0 || albumID == 0 {
		return 0, 0, errors.New("invalid track url")
	}
	return id, albumID, nil
}

func (m *Music) AddTrackToPlayList(playList yamusic.PlaylistsResult, trackURL string) {
	id, albumID, _ := m.parseTrackURL(trackURL)
	fmt.Println(id, albumID)
	tracks := []yamusic.PlaylistsTrack{{id, albumID}}

	list, _, err := m.client.Playlists().AddTracks(context.Background(), playList.Kind, playList.Revision, tracks, &yamusic.PlaylistsAddTracksOptions{At: 1})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(list.Error)
}

func (m *Music) Test() {
	uri := fmt.Sprintf("avatars/%v/download-info", 5280749)
	req, err := m.client.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	spew.Dump(req)

	//dlInfoResp := new(DownloadInfoResp)
	//resp, err := t.client.Do(ctx, req, dlInfoResp)
	//return dlInfoResp, resp, err
}

func (m *Music) GetDownloadURL(trackURL string) (string, error) {
	id, _, _ := m.parseTrackURL(trackURL)

	//m.client.Tracks().GetDownloadInfo(context.Background(), id)

	return m.client.Tracks().GetDownloadURL(context.Background(), id)
}

type AudioParams struct {
	URL          string
	Performer    string
	Title        string
	ThumbnailURL string
}

func (m *Music) GetAudioParams(trackURL string) (AudioParams, error) {
	id, _, err := m.parseTrackURL(trackURL)
	if err != nil {
		return AudioParams{}, err
	}

	trackResp, _, err := m.client.Tracks().Get(context.Background(), id)
	if err != nil {
		return AudioParams{}, err
	}
	track := trackResp.Result[0]
	//spew.Dump(track)

	params := AudioParams{}

	params.URL, err = m.client.Tracks().GetDownloadURL(context.Background(), id)
	//spew.Dump(params.URL)
	if err != nil {
		return AudioParams{}, err
	}

	for _, artist := range track.Artists {
		params.Performer += artist.Name + ", "
	}
	params.Performer = params.Performer[:len(params.Performer)-2]

	params.Title = track.Title
	params.ThumbnailURL = "https://" + track.OgImage[:len(track.OgImage)-2] + "400x400"
	//https://avatars.yandex.net/get-music-content/4384958/a4fbbb2c.a.16103149-1/320x320
	//https://avatars.yandex.net/get-music-content/8871869/46e69178.a.25384794-1/m1000x1000
	return params, nil
}
