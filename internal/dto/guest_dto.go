package dto

import "learn/internal/model"

type CreateGuestInput struct {
	Name string `json:"name" binding:"required"`
	Bio  string `json:"bio"`
}

type GuestResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	Bio  string `json:"bio"`
}

func ToGuestResponse(guest model.Guest) GuestResponse {
	return GuestResponse{
		ID:   guest.ID,
		Name: guest.Name,
		Slug: guest.Slug,
		Bio:  guest.Bio,
	}
}

func ToGuestResponses(guests []model.Guest) []GuestResponse {
	var responses []GuestResponse
	for _, guest := range guests {
		responses = append(responses, ToGuestResponse(guest))
	}
	return responses
}
