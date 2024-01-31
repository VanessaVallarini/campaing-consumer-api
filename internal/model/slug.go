package model

import (
	"time"

	"github.com/google/uuid"
)

type Slug struct {
	Id        uuid.UUID `json:"id"`
	UserId    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Active    bool      `json:"active"`
	Lat       float64   `json:"lat"`
	Long      float64   `json:"long"`
}
