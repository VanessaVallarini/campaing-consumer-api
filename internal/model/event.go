package model

import "github.com/google/uuid"

type MessageBody struct {
	Message string `json:"Message"`
}

type Event struct {
	Id          uuid.UUID `json:"id"`
	UserId      uuid.UUID `json:"user_id"`
	SlugId      uuid.UUID `json:"slug_id"`
	MerchantId  uuid.UUID `json:"merchant_id"`
	Lat         float64   `json:"lat"`
	Long        float64   `json:"long"`
	Action      string    `json:"action"`
	Active      bool      `json:"active"`
	Clicks      int       `json:"clicks"`
	Impressions int       `json:"impressions"`
}
