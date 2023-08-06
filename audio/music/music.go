package music

import (
	"context"
	"errors"
	"fmt"
	"github.com/ndrewnee/go-yamusic/yamusic"
	circuit "github.com/rubyist/circuitbreaker"
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

func (m *Music) searchTrack(id int) {
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
	matched, err := regexp.MatchString(`^https:\/\/music\.yandex\.(ru|com)\/album\/[0-9]*\/track\/[0-9]*.*`, trackURL)
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

type AudioParams struct {
	URL          string
	Performer    string
	Title        string
	ThumbnailURL string
}

//https://music.yandex.ru/album/19435876/track/95386879
//https://music.yandex.com/album/23345073/track/106849229

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
	return params, nil
}
