package usecase

import "music_bot/internal/entity"

type AudioRepo interface {
	GetFromUrl(url string) entity.Audio
}

var Audio *AudioRepo
