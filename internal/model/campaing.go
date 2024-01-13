package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type CampaingDBModel struct {
	Id          uuid.UUID `MyTag:"id"`
	UserId      uuid.UUID `MyTag:"user_id"`
	SlugId      uuid.UUID `MyTag:"slug_id"`
	MerchantId  uuid.UUID `MyTag:"merchant_id"`
	CreatedAt   time.Time `MyTag:"created_at"`
	UpdatedAt   time.Time `MyTag:"updated_at"`
	Active      bool      `MyTag:"active"`
	Lat         float64   `MyTag:"lat"`
	Long        float64   `MyTag:"long"`
	Clicks      int       `MyTag:"clicks"`
	Impressions int       `MyTag:"impressions"`
}
