package guest

import "time"

type CreateGuestInput struct {
	Name string `json:"name" binding:"required"`
	Bio  string `json:"bio"`
}

type GuestResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Bio       string    `json:"bio"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToGuestResponse(guest Guest) GuestResponse {
	return GuestResponse{
		ID:        guest.ID,
		Name:      guest.Name,
		Slug:      guest.Slug,
		Bio:       guest.Bio,
		CreatedAt: guest.CreatedAt,
		UpdatedAt: guest.UpdatedAt,
	}
}

func ToGuestResponses(guests []Guest) []GuestResponse {
	var responses []GuestResponse
	for _, guest := range guests {
		responses = append(responses, ToGuestResponse(guest))
	}
	return responses
}