package speaker

type CreateSpeakerInput struct {
	Name string `json:"name" binding:"required"`
	Bio  string `json:"bio"`
}
