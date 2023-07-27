package audio

import (
	"MusicBot/audio/music"
	"errors"
	"github.com/bogem/id3v2"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/nfnt/resize"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"sync"
)

const pwd = "" //"audio/"

type Audio struct {
	mutex  sync.Mutex
	yandex music.Music
}

func Init() *Audio {
	a := Audio{
		yandex: music.Music{},
	}
	a.yandex.Authorization()
	return &a
}

func (a *Audio) downloadFile(URL, fileName string) error {
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

func (a *Audio) deleteFileTags(fileName string) error {
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

func (a *Audio) resizePicture(fileName string) error {
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

	imgOut, err := os.Create(pwd + "newTestPicture.jpeg")
	if err != nil {
		return err
	}
	jpeg.Encode(imgOut, imgJpg, nil)
	imgOut.Close()
	return nil
}

func (a *Audio) getAudioFile(URL string) (*os.File, error) {
	err := a.downloadFile(URL, pwd+"tmpAudio.mp3")
	if err != nil {
		return nil, err
	}
	err = a.deleteFileTags(pwd + "tmpAudio.mp3")
	if err != nil {
		return nil, err
	}
	file, err := os.Open(pwd + "tmpAudio.mp3")
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (a *Audio) getPictureFile(URL string) (*os.File, error) {
	err := a.downloadFile(URL, pwd+"tmpPicture.jpeg")
	if err != nil {
		return nil, err
	}
	err = a.resizePicture(pwd + "tmpPicture.jpeg")
	if err != nil {
		return nil, err
	}
	picture, err := os.Open(pwd + "tmpPicture.jpeg")
	if err != nil {
		return nil, err
	}
	return picture, nil
}

func (a *Audio) GetParams(update *telego.Update) (*telego.SendAudioParams, error) {
	// Get chat ID from the message
	chatID := tu.ID(update.Message.Chat.ID)

	params, err := a.yandex.GetAudioParams(update.Message.Text)
	if err != nil {
		return nil, err
	}

	a.mutex.Lock()
	audioFile, err := a.getAudioFile(params.URL)
	if err != nil {
		return nil, err
	}

	pictureFile, err := a.getPictureFile(params.ThumbnailURL)
	if err != nil {
		return nil, err
	}

	audio := tu.Audio(chatID, tu.File(audioFile)).WithTitle(params.Title).WithPerformer(params.Performer)
	pictureInputFile := tu.File(pictureFile)
	audio.WithThumbnail(&pictureInputFile)
	a.mutex.Unlock()

	return audio, nil
}

/*
func (a *Audio) handlingAudioRequest(update *telego.Update) {
	audio, err := a.GetParams(update)
	if err != nil {
		a.telegram.Logger().Debugf(err.Error())
		return
	}

	_, err = a.telegram.SendAudio(audio)
	if err != nil {
		a.telegram.Logger().Debugf(err.Error())
		return
	}
}
*/
