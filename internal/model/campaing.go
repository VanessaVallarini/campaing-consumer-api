package model

import (
	"time"

	"github.com/google/uuid"
)

type Campaing struct {
	Id          uuid.UUID `json:"id"`
	UserId      uuid.UUID `json:"user_id"`
	SlugId      uuid.UUID `json:"slug_id"`
	MerchantId  uuid.UUID `json:"merchant_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Active      bool      `json:"active"`
	Lat         float64   `json:"lat"`
	Long        float64   `json:"long"`
	Clicks      int       `json:"clicks"`
	Impressions int       `json:"impressions"`
}
