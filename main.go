package main

import (
	"MusicBot/music"
	"errors"
	"fmt"
	"github.com/bogem/id3v2"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/nfnt/resize"
	"image/jpeg"
	"io"
	"net/http"
	"os"
)

//https://pkg.go.dev/github.com/ndrewnee/go-yamusic/yamusic@v1.1.0#section-sourcefiles

//func main() {
//	var m music.Music
//	m.Authorization()
//
//	m.GetAudioParams("https://music.yandex.ru/album/17678543/track/89854895")
//
//	//m.GetAudioParams("https://music.yandex.ru/album/17678543/track/89854895")
//}

//https://music.yandex.ru/album/17678543/track/89854895
//https://music.yandex.ru/album/15227716/track/81898490
//https://music.yandex.ru/album/14190852/track/78787280

func downloadFile(URL, fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}
	//Create an empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func deleteFileTags(fileName string) error {
	tag, err := id3v2.Open(fileName, id3v2.Options{Parse: true})
	if err != nil {
		return err
	}
	defer tag.Close()
	tag.DeleteAllFrames()
	if err = tag.Save(); err != nil {
		return err
	}
	return nil
}

func resizePicture(fileName string) error {
	imgIn, err := os.Open(fileName)
	if err != nil {
		return err
	}
	imgJpg, err := jpeg.Decode(imgIn)
	if err != nil {
		return err
	}
	imgIn.Close()

	imgJpg = resize.Resize(320, 320, imgJpg, resize.Bicubic)

	imgOut, err := os.Create("newTestPicture.jpeg")
	if err != nil {
		return err
	}
	jpeg.Encode(imgOut, imgJpg, nil)
	imgOut.Close()
	return nil
}

func getAudioFile(URL string) (*os.File, error) {
	err := downloadFile(URL, "tmpAudio.mp3")
	if err != nil {
		return nil, err
	}
	err = deleteFileTags("tmpAudio.mp3")
	if err != nil {
		return nil, err
	}
	file, err := os.Open("tmpAudio.mp3")
	if err != nil {
		return nil, err
	}
	return file, nil
}

func getPictureFile(URL string) (*os.File, error) {
	err := downloadFile(URL, "tmpPicture.jpeg")
	if err != nil {
		return nil, err
	}
	err = resizePicture("tmpPicture.jpeg")
	if err != nil {
		return nil, err
	}
	picture, err := os.Open("tmpPicture.jpeg")
	if err != nil {
		return nil, err
	}
	return picture, nil
}

func getSendAudioParams(chatID telego.ChatID, messageText string, m music.Music) (*telego.SendAudioParams, error) {
	params, err := m.GetAudioParams(messageText)
	if err != nil {
		return nil, err
	}

	audioFile, err := getAudioFile(params.URL)
	if err != nil {
		return nil, err
	}

	pictureFile, err := getPictureFile(params.ThumbnailURL)
	if err != nil {
		return nil, err
	}

	audio := tu.Audio(chatID, tu.File(audioFile)).WithTitle(params.Title).WithPerformer(params.Performer)
	pictureInputFile := tu.File(pictureFile)
	audio.WithThumbnail(&pictureInputFile)

	return audio, nil
}

func messageProcessing(bot *telego.Bot, update *telego.Update, m music.Music) {
	// Check if update contains a message
	if update.Message == nil {
		return
	}
	// Get chat ID from the message
	chatID := tu.ID(update.Message.Chat.ID)

	audio, err := getSendAudioParams(chatID, update.Message.Text, m)
	if err != nil {
		bot.Logger().Debugf(err.Error())
		return
	}

	_, err = bot.SendAudio(audio)
	if err != nil {
		bot.Logger().Debugf(err.Error())
		return
	}
}

func main() {
	botToken := TOKEN

	// Create Bot with debug on
	// Note: Please keep in mind that default logger may expose sensitive information, use in development only
	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get updates channel
	updates, _ := bot.UpdatesViaLongPolling(nil)

	// Stop reviving updates from update channel
	defer bot.StopLongPolling()

	var m music.Music
	m.Authorization()

	//loger := bot.Logger()
	// Loop through all updates when they came
	for update := range updates {
		messageProcessing(bot, &update, m)
	}
}

//
/*
type AudioParams struct {
	URL          string
	Caption      string
	Performer    string
	Title        string
	ThumbnailURL string
}
*/
