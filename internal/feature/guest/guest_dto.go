package guest

type CreateGuestInput struct {
	Name string `json:"name" binding:"required"`
	Bio  string `json:"bio"`
}