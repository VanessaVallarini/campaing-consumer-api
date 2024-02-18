package repository

import (
	"campaing-comsumer-service/internal/metrics"
	"campaing-comsumer-service/internal/model"
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lockp111/go-easyzap"
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
		easyzap.Warn("select merchant %v cancel by context. msg: %v", param, err)
		mv := []string{"GetById", "error", "starts_transaction"}
		c.metrics.MerchantRepository.WithLabelValues(mv...).Inc()

		return model.Merchant{}, err
	}

	var merchant model.Merchant
	row := tx.QueryRowContext(ctx, "select * from merchant where id = $1;", param)

	err = row.Scan(&merchant.Id, &merchant.UserId, &merchant.SlugId, &merchant.CreatedAt, &merchant.UpdatedAt, &merchant.Name, &merchant.Active, &merchant.Lat, &merchant.Long)
	if err != nil {
		if err == sql.ErrNoRows {
			easyzap.Warn("no lines found merchant for merchant_id: %v", param, err)
			tx.Rollback()
			mv := []string{"GetById", "error", "no_rows"}
			c.metrics.MerchantRepository.WithLabelValues(mv...).Inc()

			return model.Merchant{}, nil
		}
		easyzap.Errorf("scan merchant %v fail. msg: %v", param, err)
		tx.Rollback()
		mv := []string{"GetById", "error", "scan"}
		c.metrics.MerchantRepository.WithLabelValues(mv...).Inc()

		return model.Merchant{}, err
	}

	err = tx.Commit()
	if err != nil {
		easyzap.Warn("select merchant id %v fail. msg: %v", param, err)
		mv := []string{"GetById", "error", "commit"}
		c.metrics.MerchantRepository.WithLabelValues(mv...).Inc()

		return model.Merchant{}, err
	}

	mv := []string{"GetById", "success", ""}
	c.metrics.MerchantRepository.WithLabelValues(mv...).Inc()

	return merchant, nil
}
