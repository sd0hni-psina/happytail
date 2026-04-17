package models

import (
	"time"
)

type ListingType string
type PostStatus string

const (
	ListingTypeSale ListingType = "sale"
	ListingTypeGive ListingType = "give"

	PostStatusActive   PostStatus = "active"
	PostStatusInactive PostStatus = "inactive"
	PostStatusDeleted  PostStatus = "deleted"
)

type Money struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type Post struct {
	ID          int         `json:"id"`
	UserID      int         `json:"user_id"`
	AnimalID    int         `json:"animal_id"`
	ListingType ListingType `json:"listing_type"`
	Price       *Money      `json:"price"`
	Reason      *string     `json:"reason"`
	PhotoURLs   []string    `json:"photo_urls"`
	ContactInfo string      `json:"contact_info"`
	Status      PostStatus  `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type PostInput struct {
	UserID      int
	AnimalID    int
	ListingType ListingType
	Price       *Money
	Reason      *string
	PhotoURLs   []string
	ContactInfo string
}
