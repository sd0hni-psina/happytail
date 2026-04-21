package models

import (
	"errors"
	"fmt"
	"strings"
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

var validListingType = map[ListingType]struct{}{
	ListingTypeGive: {},
	ListingTypeSale: {},
}

func (lt ListingType) Normalize() ListingType {
	return ListingType(strings.ToLower(string(lt)))
}

func (lt ListingType) IsValid() bool {
	_, ok := validListingType[lt]
	return ok
}

type Money struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type Post struct {
	ID          int         `json:"-"`
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

type CreatePostInput struct {
	UserID      int         `json:"-"`
	AnimalID    int         `json:"animal_id"`
	ListingType ListingType `json:"listing_type"`
	Price       *Money      `json:"price"`
	Reason      *string     `json:"reason"`
	PhotoURLs   []string    `json:"photo_urls"`
	ContactInfo string      `json:"contact_info"`
}

func (cpi *CreatePostInput) Validate() error {
	cpi.ListingType = cpi.ListingType.Normalize()
	validationErrors := make(map[string]string)

	if cpi.AnimalID == 0 {
		validationErrors["Animal ID"] = "Animal ID is required, connot be 0"
	}
	if cpi.ListingType == "" {
		validationErrors["Listing Type"] = "Listing type cannot be empty. For sale or give, only!"
	} else if !cpi.ListingType.IsValid() {
		validationErrors["Listing Type"] = fmt.Sprintf("invalid type: %s", cpi.ListingType)
	}
	if cpi.ContactInfo == "" {
		validationErrors["Contact Info"] = "Contact info is required"
	}
	if cpi.ListingType == ListingTypeSale && cpi.Price == nil {
		validationErrors["Price"] = "Price is required if Listing Type for sale"
	}
	if len(validationErrors) > 0 {
		msgs := make([]string, 0, len(validationErrors))
		for field, msg := range validationErrors {
			msgs = append(msgs, field+": "+msg)
		}
		return errors.New(strings.Join(msgs, ", "))
	}
	return nil
}
