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

type Slug struct {
	db      *pgxpool.Pool
	metrics *metrics.Metrics
}

func NewSlugRepository(metrics *metrics.Metrics, db *pgxpool.Pool) *Slug {
	return &Slug{
		db:      db,
		metrics: metrics,
	}
}

func (c *Slug) GetById(ctx context.Context, param uuid.UUID) (model.Slug, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	var slug model.Slug
	rows := c.db.QueryRow(ctx, "select * from slug where id = $1;", param)

	err := rows.Scan(&slug.Id, &slug.UserId, &slug.CreatedAt, &slug.UpdatedAt, &slug.Active, &slug.Lat, &slug.Long)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return model.Slug{}, nil
		}
		mv := []string{"GetById", "error", "scan"}
		c.metrics.SlugRepository.WithLabelValues(mv...).Inc()
		errWrap := errors.Wrapf(err, "scan slug %v fail", param)
		easyzap.Errorf("scan slug %v fail. msg: %v", param, errWrap)
		return model.Slug{}, errWrap
	}

	mv := []string{"GetById", "success", ""}
	c.metrics.SlugRepository.WithLabelValues(mv...).Inc()

	return slug, nil
}
