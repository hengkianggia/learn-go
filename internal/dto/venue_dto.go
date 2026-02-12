package dto

import "learn/internal/model"

type CreateVenueInput struct {
	Name     string `json:"name" binding:"required"`
	Address  string `json:"address" binding:"required"`
	City     string `json:"city"`
	State    string `json:"state"`
	ZipCode  string `json:"zip_code"`
	Capacity int    `json:"capacity"`
	IsActive bool   `json:"is_active,omitempty"`
	Country  string `json:"country,omitempty"`
}

type UpdateVenueInput struct {
	Name     *string `json:"name,omitempty"`
	Address  *string `json:"address,omitempty"`
	City     *string `json:"city,omitempty"`
	State    *string `json:"state,omitempty"`
	ZipCode  *string `json:"zip_code,omitempty"`
	Capacity *int    `json:"capacity,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
	Country  *string `json:"country,omitempty"`
}

type VenueResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Address  string `json:"address"`
	City     string `json:"city"`
	State    string `json:"state"`
	ZipCode  string `json:"zip_code"`
	Capacity int    `json:"capacity"`
	IsActive bool   `json:"is_active"`
	Country  string `json:"country"`
}

func ToVenueResponse(venue model.Venue) VenueResponse {
	return VenueResponse{
		ID:       venue.ID,
		Name:     venue.Name,
		Slug:     venue.Slug,
		Address:  venue.Address,
		City:     venue.City,
		State:    venue.State,
		ZipCode:  venue.ZipCode,
		Capacity: venue.Capacity,
		IsActive: venue.IsActive,
		Country:  venue.Country,
	}
}

func ToVenueResponses(venues []model.Venue) []VenueResponse {
	var responses []VenueResponse
	for _, venue := range venues {
		responses = append(responses, ToVenueResponse(venue))
	}
	return responses
}
