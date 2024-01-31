package repository

import (
	"campaing-comsumer-service/internal/model"
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lockp111/go-easyzap"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type Slug struct {
	conn *sql.DB
}

func NewSlugRepository(conn *sql.DB) *Slug {
	return &Slug{
		conn: conn,
	}
}

func (c *Slug) GetById(param uuid.UUID) (model.Slug, error) {
	span, ctx := tracer.StartSpanFromContext(context.Background(), "slug_repository.get-by-id",
		tracer.ResourceName("postgres"),
		tracer.SpanType("db"),
	)
	defer span.Finish()

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	tx, err := c.conn.BeginTx(ctx, nil)
	if err != nil {
		easyzap.Errorf("select slug %v cancel by context. msg: %v", param, err)
		span.Finish(tracer.WithError(err))

		return model.Slug{}, err
	}

	var slug model.Slug
	row := tx.QueryRowContext(ctx, "select * from slug where id = $1;", param)

	err = row.Scan(&slug.Id, &slug.UserId, &slug.CreatedAt, &slug.UpdatedAt, &slug.Active, &slug.Lat, &slug.Long)
	if err != nil {
		if err == sql.ErrNoRows {
			easyzap.Warn("no lines found for slug_id: %v", param, err)
			tx.Rollback()

			return model.Slug{}, nil

		}
		easyzap.Errorf("scan slug %v fail. msg: %v", param, err)
		tx.Rollback()
		span.Finish(tracer.WithError(err))

		return model.Slug{}, err

	}

	err = tx.Commit()
	if err != nil {
		easyzap.Errorf("select slug %v fail. msg: %v", param, err)
		span.Finish(tracer.WithError(err))

		return model.Slug{}, err
	}

	return slug, nil
}
