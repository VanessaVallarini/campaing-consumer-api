package repository

import (
	"campaing-comsumer-service/internal/metrics"
	"campaing-comsumer-service/internal/model"
	"context"
	"time"

	"github.com/google/uuid"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lockp111/go-easyzap"
	"github.com/pkg/errors"
)

type Merchant struct {
	db      *pgxpool.Pool
	metrics *metrics.Metrics
}

func NewMerchantRepository(metrics *metrics.Metrics, db *pgxpool.Pool) *Merchant {
	return &Merchant{
		db:      db,
		metrics: metrics,
	}
}

func (c *Merchant) GetById(ctx context.Context, param uuid.UUID) (model.Merchant, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	var merchant model.Merchant
	rows := c.db.QueryRow(ctx, "select * from merchant where id = $1;", param)

	err := rows.Scan(&merchant.Id, &merchant.UserId, &merchant.SlugId, &merchant.CreatedAt, &merchant.UpdatedAt, &merchant.Name, &merchant.Active, &merchant.Lat, &merchant.Long)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return model.Merchant{}, nil
		}
		mv := []string{"GetById", "error", "scan"}
		c.metrics.MerchantRepository.WithLabelValues(mv...).Inc()
		errWrap := errors.Wrapf(err, "scan merchant %v fail", param)
		easyzap.Errorf("scan merchant %v fail. msg: %v", param, errWrap)
		return model.Merchant{}, errWrap
	}

	mv := []string{"GetById", "success", ""}
	c.metrics.MerchantRepository.WithLabelValues(mv...).Inc()

	return merchant, nil
}
