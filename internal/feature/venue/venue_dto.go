package venue

type CreateVenueInput struct {
	Name      string `json:"name" binding:"required"`
	Address   string `json:"address" binding:"required"`
	City      string `json:"city"`
	State     string `json:"state"`
	ZipCode   string `json:"zip_code"`
	Capacity  int    `json:"capacity"`
	IsActive  bool   `json:"is_active,omitempty"`
}
