package models

import advertisements "api-tests-template/internal/managers/advertisements/models"

type AddRequest struct {
	AdvertisementID string `json:"advertisement_id"`
}

type Favorite struct {
	FavoriteID    string                       `json:"favorite_id"`
	AddedAt       string                       `json:"added_at"`
	Advertisement advertisements.Advertisement `json:"advertisement"`
}

type List struct {
	Items      []Favorite `json:"items"`
	Page       int        `json:"page"`
	PageSize   int        `json:"page_size"`
	Total      int        `json:"total"`
	TotalPages int        `json:"total_pages"`
}
