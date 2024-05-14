package usecase

import (
	"music_bot/internal/entity"
)

type Sender interface {
	UnidentSecretMsg(id entity.UserID)
	ChildRegGreeting(id entity.UserID)
	HostRegGreeting(id entity.UserID)
	SendStartMenu(id entity.UserID)
	SendMenu(id entity.UserID)
	Out(id entity.UserID)
	UserConnect(hostId entity.UserID, childId entity.UserID)
	UserDisconnect(hostId entity.UserID, childId entity.UserID)
	TrackNotFound(id entity.UserID)
	TrackAddToPlayList(id entity.UserID, audio entity.Audio)
	OptionInDevelopment(id entity.UserID)
	SendInfo(id entity.UserID)
	UnknownError(id entity.UserID)
}
