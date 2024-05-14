package usecase

import (
  "music_bot/internal/entity"
)

type Sender interface {
  SendMenu(id entity.UserID)
  Out(id entity.UserID)
  UserConnect(hostId entity.UserID, childId entity.UserID)
  UserDisconnect(hostId entity.UserID, childId entity.UserID)
  TrackNotFound(id entity.UserID)
  TrackAddToPlayList(id entity.UserID, audio entity.Audio)
  OptionInDevelopment(id entity.UserID)
  SendInfo(id entity.UserID)
} 