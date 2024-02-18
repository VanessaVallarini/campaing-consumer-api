package repository

import (
	"campaing-comsumer-service/internal/metrics"
	"campaing-comsumer-service/internal/model"
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lockp111/go-easyzap"
	"github.com/pkg/errors"
)

type Merchant struct {
	conn    *sql.DB
	metrics *metrics.Metrics
}

func NewMerchantRepository(metrics *metrics.Metrics, conn *sql.DB) *Merchant {
	return &Merchant{
		conn:    conn,
		metrics: metrics,
	}
}

func (c *Merchant) GetById(ctx context.Context, param uuid.UUID) (model.Merchant, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	tx, err := c.conn.BeginTx(ctx, nil)
	if err != nil {
		mv := []string{"GetById", "error", "starts_transaction"}
		c.metrics.MerchantRepository.WithLabelValues(mv...).Inc()
		errWrap := errors.Wrapf(err, "select merchant %v cancel by context", param)
		easyzap.Warn(ctx, errWrap, "select merchant cancel by context")
		return model.Merchant{}, errWrap
	}

	var merchant model.Merchant
	row := tx.QueryRowContext(ctx, "select * from merchant where id = $1;", param)

	err = row.Scan(&merchant.Id, &merchant.UserId, &merchant.SlugId, &merchant.CreatedAt, &merchant.UpdatedAt, &merchant.Name, &merchant.Active, &merchant.Lat, &merchant.Long)
	if err != nil {
		if err == sql.ErrNoRows {
			tx.Rollback()
			return model.Merchant{}, nil
		}
		tx.Rollback()
		mv := []string{"GetById", "error", "scan"}
		c.metrics.MerchantRepository.WithLabelValues(mv...).Inc()
		errWrap := errors.Wrapf(err, "scan merchant %v fail", param)
		easyzap.Error(ctx, errWrap, "scan merchant fail")
		return model.Merchant{}, errWrap
	}

	err = tx.Commit()
	if err != nil {
		mv := []string{"GetById", "error", "commit"}
		c.metrics.MerchantRepository.WithLabelValues(mv...).Inc()
		errWrap := errors.Wrapf(err, "select merchant id %v fail", param)
		easyzap.Error(ctx, errWrap, "select merchant fail")
		return model.Merchant{}, errWrap
	}

	mv := []string{"GetById", "success", ""}
	c.metrics.MerchantRepository.WithLabelValues(mv...).Inc()

	return merchant, nil
}
