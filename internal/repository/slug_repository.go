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
		easyzap.Errorf("select slug %v cancel by context. msg: %v", param, err)
		mv := []string{"GetById", "error", "starts_transaction"}
		c.metrics.SlugRepository.WithLabelValues(mv...).Inc()

		return model.Slug{}, err
	}

	var slug model.Slug
	row := tx.QueryRowContext(ctx, "select * from slug where id = $1;", param)

	err = row.Scan(&slug.Id, &slug.UserId, &slug.CreatedAt, &slug.UpdatedAt, &slug.Active, &slug.Lat, &slug.Long)
	if err != nil {
		if err == sql.ErrNoRows {
			easyzap.Warn("no lines found for slug_id: %v", param, err)
			tx.Rollback()
			mv := []string{"GetById", "error", "no_rows"}
			c.metrics.SlugRepository.WithLabelValues(mv...).Inc()

			return model.Slug{}, nil

		}
		easyzap.Errorf("scan slug %v fail. msg: %v", param, err)
		tx.Rollback()
		mv := []string{"GetById", "error", "scan"}
		c.metrics.SlugRepository.WithLabelValues(mv...).Inc()

		return model.Slug{}, err

	}

	err = tx.Commit()
	if err != nil {
		easyzap.Errorf("select slug %v fail. msg: %v", param, err)
		mv := []string{"GetById", "error", "commit"}
		c.metrics.SlugRepository.WithLabelValues(mv...).Inc()

		return model.Slug{}, err
	}

	mv := []string{"GetById", "success", ""}
	c.metrics.SlugRepository.WithLabelValues(mv...).Inc()

	return slug, nil
}
