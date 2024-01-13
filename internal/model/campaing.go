package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type Campaing struct {
	Id          uuid.UUID
	UserId      uuid.UUID
	SlugId      uuid.UUID
	MerchantId  uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Active      bool
	Lat         float64
	Long        float64
	Clicks      int
	Impressions int
}
