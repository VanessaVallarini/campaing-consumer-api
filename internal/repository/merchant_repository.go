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

type Merchant struct {
	conn *sql.DB
}

func NewMerchantRepository(conn *sql.DB) *Merchant {
	return &Merchant{
		conn: conn,
	}
}

func (c *Merchant) GetById(param uuid.UUID) (model.Merchant, error) {
	span, ctx := tracer.StartSpanFromContext(context.Background(), "merchant_repository.get-by-id",
		tracer.ResourceName("postgres"),
		tracer.SpanType("db"),
	)
	defer span.Finish()

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	tx, err := c.conn.BeginTx(ctx, nil)
	if err != nil {
		easyzap.Warn("select merchant %v cancel by context. msg: %v", param, err)
		span.Finish(tracer.WithError(err))

		return model.Merchant{}, err
	}

	var merchant model.Merchant
	row := tx.QueryRowContext(ctx, "select * from merchant where id = $1;", param)

	err = row.Scan(&merchant.Id, &merchant.UserId, &merchant.SlugId, &merchant.CreatedAt, &merchant.UpdatedAt, &merchant.Name, &merchant.Active, &merchant.Lat, &merchant.Long)
	if err != nil {
		if err == sql.ErrNoRows {
			easyzap.Warn("no lines found merchant for merchant_id: %v", param, err)
			tx.Rollback()

			return model.Merchant{}, nil
		}
		easyzap.Errorf("scan merchant %v fail. msg: %v", param, err)
		tx.Rollback()
		span.Finish(tracer.WithError(err))

		return model.Merchant{}, err
	}

	err = tx.Commit()
	if err != nil {
		easyzap.Warn("select merchant id %v fail. msg: %v", param, err)
		span.Finish(tracer.WithError(err))

		return model.Merchant{}, err
	}

	return merchant, nil
}
