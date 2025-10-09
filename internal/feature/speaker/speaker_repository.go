package speaker

import "gorm.io/gorm"

type SpeakerRepository interface {
	CreateSpeaker(speaker *Speaker) error
	GetSpeakerByID(id uint) (*Speaker, error)
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

func (r *speakerRepository) GetSpeakerByID(id uint) (*Speaker, error) {
	var speaker Speaker
	err := r.db.First(&speaker, id).Error
	return &speaker, err
}