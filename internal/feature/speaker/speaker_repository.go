package speaker

import "gorm.io/gorm"

type SpeakerRepository interface {
	CreateSpeaker(speaker *Speaker) error
}

type speakerRepository struct {
	db *gorm.DB
}

func NewSpeakerRepository(db *gorm.DB) SpeakerRepository {
	return &speakerRepository{db: db}
}

func (r *speakerRepository) CreateSpeaker(speaker *Speaker) error {
	return r.db.Create(speaker).Error
}
