package dto

import "learn/internal/model"

// input
type CreateGuestInput struct {
	Name string `json:"name" binding:"required"`
	Bio  string `json:"bio"`
}

type UpdateGuestInput struct {
	Name *string `json:"name,omitempty"`
	Bio  *string `json:"bio,omitempty"`
}

// response
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

// simple response
type GuestSimpleResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func ToGuestSimpleResponse(guest model.Guest) GuestSimpleResponse {
	return GuestSimpleResponse{
		ID:   guest.ID,
		Name: guest.Name,
		Slug: guest.Slug,
	}
}

func ToGuestSimpleResponses(guests []model.Guest) []GuestSimpleResponse {
	var responses []GuestSimpleResponse
	for _, guest := range guests {
		responses = append(responses, ToGuestSimpleResponse(guest))
	}
	return responses
}
