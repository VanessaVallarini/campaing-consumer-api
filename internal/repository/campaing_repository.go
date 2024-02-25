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

type Campaing struct {
	db      *pgxpool.Pool
	metrics *metrics.Metrics
}

func NewCampaingRepository(metrics *metrics.Metrics, db *pgxpool.Pool) *Campaing {
	return &Campaing{
		db:      db,
		metrics: metrics,
	}
}

func (c *Campaing) GetByMerchantId(ctx context.Context, param uuid.UUID) (model.Campaing, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	var campaing model.Campaing
	rows := c.db.QueryRow(ctx, "select * from campaing where merchant_id = $1 and active = true;", param)

	err := rows.Scan(&campaing.Id, &campaing.UserId, &campaing.SlugId, &campaing.MerchantId, &campaing.CreatedAt, &campaing.UpdatedAt, &campaing.Active, &campaing.Lat, &campaing.Long, &campaing.Clicks, &campaing.Impressions)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return model.Campaing{}, nil
		}
		mv := []string{"GetByMerchantId", "error", "scan"}
		c.metrics.CampaingRepository.WithLabelValues(mv...).Inc()
		errWrap := errors.Wrapf(err, "scan campaing %v fail", param)
		easyzap.Errorf("scan campaing by merchatn id %v fail. msg: %v", campaing.MerchantId, errWrap)
		return model.Campaing{}, errWrap
	}

	mv := []string{"GetByMerchantId", "success", ""}
	c.metrics.CampaingRepository.WithLabelValues(mv...).Inc()

	return campaing, nil
}

func (c *Campaing) GetById(ctx context.Context, param uuid.UUID) (model.Campaing, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	var campaing model.Campaing
	rows := c.db.QueryRow(ctx, "select * from campaing where id = $1 and active = true;", param)

	err := rows.Scan(&campaing.Id, &campaing.UserId, &campaing.SlugId, &campaing.MerchantId, &campaing.CreatedAt, &campaing.UpdatedAt, &campaing.Active, &campaing.Lat, &campaing.Long, &campaing.Clicks, &campaing.Impressions)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return model.Campaing{}, nil
		}
		mv := []string{"GetById", "error", "scan"}
		c.metrics.CampaingRepository.WithLabelValues(mv...).Inc()
		errWrap := errors.Wrapf(err, "scan campaing %v fail", param)
		easyzap.Errorf("scan campaing by id %v fail. msg: %v", campaing.MerchantId, errWrap)
		return model.Campaing{}, errWrap
	}

	mv := []string{"GetById", "success", ""}
	c.metrics.CampaingRepository.WithLabelValues(mv...).Inc()

	return campaing, nil
}

func (c *Campaing) Create(ctx context.Context, tx pgx.Tx, params model.Campaing) error {
	_, err := tx.Exec(ctx, "insert into campaing(id,user_id,slug_id,merchant_id,created_at,updated_at,active,lat,long,clicks,impressions) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)",
		params.Id, params.UserId, params.SlugId, params.MerchantId, params.CreatedAt, params.UpdatedAt, params.Active, params.Lat, params.Long, params.Clicks, params.Impressions)
	if err != nil {
		mv := []string{"Create", "error", "exec"}
		c.metrics.CampaingRepository.WithLabelValues(mv...).Inc()
		errWrap := errors.Wrapf(err, "create campaing %v query exec fail", params)
		easyzap.Errorf("create campaing %v query exec fail. msg: %v", params.MerchantId, errWrap)
		return errWrap
	}

	mv := []string{"Create", "success", ""}
	c.metrics.CampaingRepository.WithLabelValues(mv...).Inc()

	return nil
}

func (c *Campaing) Update(ctx context.Context, tx pgx.Tx, params model.Campaing) error {
	_, err := tx.Exec(ctx, "update campaing set user_id=$2,slug_id=$3,updated_at=$4,active=$5,lat=$6,long=$7,clicks=$8,impressions=$9 where id = $1",
		params.Id, params.UserId, params.SlugId, params.UpdatedAt, params.Active, params.Lat, params.Long, params.Clicks, params.Impressions)
	if err != nil {
		mv := []string{"Update", "error", "exec"}
		c.metrics.CampaingRepository.WithLabelValues(mv...).Inc()
		errWrap := errors.Wrapf(err, "update campaing %v query exec fail", params)
		easyzap.Errorf("update campaing %v query exec fail. msg: %v", params.MerchantId, errWrap)
		return errWrap
	}

	mv := []string{"Update", "success", ""}
	c.metrics.CampaingRepository.WithLabelValues(mv...).Inc()

	return nil
}

func (c *Campaing) Delete(ctx context.Context, tx pgx.Tx, param uuid.UUID) error {
	_, err := tx.Exec(ctx, "delete from campaing where id = $1", param)
	if err != nil {
		mv := []string{"Delete", "error", "exec"}
		c.metrics.CampaingRepository.WithLabelValues(mv...).Inc()
		errWrap := errors.Wrapf(err, "delete campaing id %v fail", param)
		easyzap.Errorf("delete campaing %v query exec fail. msg: %v", param, errWrap)
		return errWrap
	}

	mv := []string{"Delete", "success", ""}
	c.metrics.CampaingRepository.WithLabelValues(mv...).Inc()
	return nil
}
