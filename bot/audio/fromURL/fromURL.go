package fromURL

import (
	"MusicBot/log"
	"MusicBot/passGen"
	"MusicBot/telegram"
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

const pwd = "audio/tmp/"

func downloadFile(URL, fileName string) error {
	//GetMsg the response bytes from the url
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

	imgOut, err := os.Create(fileName)
	if err != nil {
		return err
	}
	jpeg.Encode(imgOut, imgJpg, nil)
	imgOut.Close()
	return nil
}

func getAudioFile(audioInf AudioInf) (*os.File, error) {
	fileName := fmt.Sprintf("%v%v;%v;%v.mp3",
		pwd,
		audioInf.Title(),
		audioInf.Performer(),
		passGen.GeneratePassword(3, 0, 0, 0),
	)
	err := downloadFile(audioInf.URL(), fileName)
	if err != nil {
		return nil, err
	}
	err = deleteFileTags(fileName)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func getPictureFile(audioInf AudioInf) (*os.File, error) {
	fileName := fmt.Sprintf("%v%v;%v;%v.jpeg",
		pwd,
		audioInf.Title(),
		audioInf.Performer(),
		passGen.GeneratePassword(3, 0, 0, 0),
	)
	err := downloadFile(audioInf.ThumbnailURL(), fileName)
	if err != nil {
		return nil, err
	}
	err = resizePicture(fileName)
	if err != nil {
		return nil, err
	}
	picture, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	return picture, nil
}

type AudioInf interface {
	URL() string
	Title() string
	Performer() string
	ThumbnailURL() string
}

func deleteFile(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Logger.Error(err)
	}
	err = os.Remove(file.Name())
	if err != nil {
		log.Logger.Error(err)
	}
}

func FromURL(audioInf AudioInf) (audio *telego.SendAudioParams, err error) {
	chatID := tu.ID(-809440484)

	audioFile, err := getAudioFile(audioInf)
	if err != nil {
		return nil, err
	}
	defer deleteFile(audioFile)

	audio = tu.Audio(chatID,
		tu.File(audioFile)).
		WithTitle(audioInf.Title()).
		WithPerformer(audioInf.Performer())

	pictureFile, err := getPictureFile(audioInf)
	if err == nil {
		defer deleteFile(pictureFile)
		pictureInputFile := tu.File(pictureFile)
		audio.WithThumbnail(&pictureInputFile)
	}

	sendAudio, err := telegram.TG.SendAudio(audio)
	if err != nil {
		return nil, err
	}

	audio.Audio.File = nil
	audio.Thumbnail = nil
	audio.Audio.FileID = sendAudio.Audio.FileID

	return audio, nil
}
