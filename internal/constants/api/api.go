package api

import "time"

const (
	RequestTimeout       = 30 * time.Second
	MaxTitleLength       = 50
	MaxDescriptionLength = 500
	MaxPrice             = 1000000
	MaxQuantity          = 100
	MaxPhotos            = 3
	MaxFavorites         = 50
	PageSize             = 50
	InvalidToken         = "incorrect_token"
	InvalidID            = "not-a-uuid"

	InvalidRequest = "invalid_request"
	Unauthorized   = "unauthorized"
	Forbidden      = "forbidden"
	NotFound       = "not_found"
	AlreadyExists  = "already_exists"
	LimitReached   = "limit_reached"

	AuthorizationRequired       = "Authorization header required"
	InvalidAuthorizationMessage = "Invalid or expired token"
	AdvertisementNotFound       = "Advertisement not found"
	FavoriteNotFound            = "Favorite entry not found"
	OwnAdvertisement            = "Cannot add own advertisement to favorites"
	DuplicateFavorite           = "Advertisement is already in favorites"
	FavoritesLimitMessage       = "Favorites limit of 50 reached"
	TitleRequired               = "title is required"
	DescriptionRequired         = "description is required"
	PriceOutOfRange             = "price must be between 0 and 1000000 inclusive"
	QuantityOutOfRange          = "quantity must be between 1 and 100 inclusive"
)
