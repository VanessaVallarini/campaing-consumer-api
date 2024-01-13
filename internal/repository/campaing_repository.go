package repository

import (
	"campaing-comsumer-service/internal/db"
)

type CampaingRepository struct {
	conn db.IDb
}

func NewFormRepository(conn db.IDb) *CampaingRepository {
	return &CampaingRepository{
		conn: conn,
	}
}
