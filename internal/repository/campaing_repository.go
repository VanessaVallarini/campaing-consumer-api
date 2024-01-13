package repository

import (
	"campaing-comsumer-service/internal/db"
	"campaing-comsumer-service/internal/model"
	"context"
	"fmt"

	"github.com/lockp111/go-easyzap"
	"go.uber.org/zap"
)

const TABLE_NAME = "campaing"

type ICampaingRepository interface {
	Create(ctx context.Context, params model.Campaing) error
}

type CampaingRepository struct {
	conn db.IDb
}

func NewFormRepository(conn db.IDb) *CampaingRepository {
	return &CampaingRepository{
		conn: conn,
	}
}

func (repo *CampaingRepository) Create(ctx context.Context, params model.Campaing) error {
	result, err := repo.conn.Exec(ctx, "insert into campaing(id,user_id,slug_id,merchant_id,active,lat,long) VALUES($1,$2,$3,$4,$5,$6,$7)",
		params.Id, params.UserId, params.SlugId, params.MerchantId, params.CreatedAt, params.UpdatedAt, params.Active, params.Lat, params.Long)
	if err != nil {
		id, err := result.LastInsertId()
		easyzap.Error(ctx, err, "error on create query for campaing id",
			zap.String("campaingId", fmt.Sprintf("%v", id)))
	}

	return nil
}
