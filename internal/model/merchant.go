package model

import (
	"time"

	"github.com/google/uuid"
)

type Merchant struct {
	Id        uuid.UUID `json:"id"`
	UserId    uuid.UUID `json:"user_id"`
	SlugId    uuid.UUID `json:"slug_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Active    bool      `json:"active"`
	Lat       float64   `json:"lat"`
	Long      float64   `json:"long"`
}
