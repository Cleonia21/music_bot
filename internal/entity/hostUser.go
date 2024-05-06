package entity

import "music_bot/utils/queueCollection"

type HostUser struct {
	ID       UserID
	playList queueCollection.QueueCollections[UserID, Audio]
	pass     string
}

func NewHostUser(id UserID) (hu *HostUser) {
	return hu
}

func (u *HostUser) SetAudio(audio Audio) {

}

func (u *HostUser) GetAudio() (audio []Audio) {
	return audio
}

func (u *HostUser) GetAudioNum(id UserID) {

}

func (u *HostUser) ValidatePass(pass string) (ok bool) {
	return true
}
