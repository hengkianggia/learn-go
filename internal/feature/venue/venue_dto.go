package venue

import "time"

type CreateVenueInput struct {
	Name      string `json:"name" binding:"required"`
	Address   string `json:"address" binding:"required"`
	City      string `json:"city"`
	State     string `json:"state"`
	ZipCode   string `json:"zip_code"`
	Capacity  int    `json:"capacity"`
	IsActive  bool   `json:"is_active,omitempty"`
}

type VenueResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Address   string    `json:"address"`
	City      string    `json:"city"`
	State     string    `json:"state"`
	ZipCode   string    `json:"zip_code"`
	Capacity  int       `json:"capacity"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToVenueResponse(venue Venue) VenueResponse {
	return VenueResponse{
		ID:        venue.ID,
		Name:      venue.Name,
		Slug:      venue.Slug,
		Address:   venue.Address,
		City:      venue.City,
		State:     venue.State,
		ZipCode:   venue.ZipCode,
		Capacity:  venue.Capacity,
		IsActive:  venue.IsActive,
		CreatedAt: venue.CreatedAt,
		UpdatedAt: venue.UpdatedAt,
	}
}

func ToVenueResponses(venues []Venue) []VenueResponse {
	var responses []VenueResponse
	for _, venue := range venues {
		responses = append(responses, ToVenueResponse(venue))
	}
	return responses
}