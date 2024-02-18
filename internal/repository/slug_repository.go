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

type Slug struct {
	conn    *sql.DB
	metrics *metrics.Metrics
}

func NewSlugRepository(metrics *metrics.Metrics, conn *sql.DB) *Slug {
	return &Slug{
		conn:    conn,
		metrics: metrics,
	}
}

func (c *Slug) GetById(ctx context.Context, param uuid.UUID) (model.Slug, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	tx, err := c.conn.BeginTx(ctx, nil)
	if err != nil {
		mv := []string{"GetById", "error", "starts_transaction"}
		c.metrics.SlugRepository.WithLabelValues(mv...).Inc()
		errWrap := errors.Wrapf(err, "select slug %v cancel by context", param)
		easyzap.Warn(ctx, errWrap, "select slug cancel by context")
		return model.Slug{}, errWrap
	}

	var slug model.Slug
	row := tx.QueryRowContext(ctx, "select * from slug where id = $1;", param)

	err = row.Scan(&slug.Id, &slug.UserId, &slug.CreatedAt, &slug.UpdatedAt, &slug.Active, &slug.Lat, &slug.Long)
	if err != nil {
		if err == sql.ErrNoRows {
			tx.Rollback()
			return model.Slug{}, nil
		}
		tx.Rollback()
		mv := []string{"GetById", "error", "scan"}
		c.metrics.SlugRepository.WithLabelValues(mv...).Inc()
		errWrap := errors.Wrapf(err, "scan slug %v fail", param)
		easyzap.Error(ctx, errWrap, "scan slug fail")
		return model.Slug{}, errWrap
	}

	err = tx.Commit()
	if err != nil {
		mv := []string{"GetById", "error", "commit"}
		c.metrics.SlugRepository.WithLabelValues(mv...).Inc()
		errWrap := errors.Wrapf(err, "select slug id %v fail", param)
		easyzap.Error(ctx, errWrap, "select slug fail")
		return model.Slug{}, errWrap
	}

	mv := []string{"GetById", "success", ""}
	c.metrics.SlugRepository.WithLabelValues(mv...).Inc()

	return slug, nil
}
