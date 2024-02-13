package repository

import (
	"campaing-comsumer-service/internal/model"
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lockp111/go-easyzap"
)

type Campaing struct {
	conn *sql.DB
}

func NewCampaingRepository(conn *sql.DB) *Campaing {
	return &Campaing{
		conn: conn,
	}
}

func (c *Campaing) GetByMerchantId(ctx context.Context, param uuid.UUID) (model.Campaing, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	tx, err := c.conn.BeginTx(ctx, nil)
	if err != nil {
		easyzap.Warn("select campaing by merchant_id %v cancel by context. msg: %v", param, err)

		return model.Campaing{}, err
	}

	var campaing model.Campaing
	row := tx.QueryRowContext(ctx, "select * from campaing where merchant_id = $1 and active = true;", param)

	err = row.Scan(&campaing.Id, &campaing.UserId, &campaing.SlugId, &campaing.MerchantId, &campaing.CreatedAt, &campaing.UpdatedAt, &campaing.Active, &campaing.Lat, &campaing.Long, &campaing.Clicks, &campaing.Impressions)
	if err != nil {
		if err == sql.ErrNoRows {
			easyzap.Warn("no lines found campaing for merchant_id: %v", param, err)
			tx.Rollback()

			return model.Campaing{}, nil
		}
		easyzap.Errorf("scan campaing %v fail. msg: %v", param, err)
		tx.Rollback()

		return model.Campaing{}, err
	}

	err = tx.Commit()
	if err != nil {
		easyzap.Warn("select campaing by merchant_id %v fail. msg: %v", param, err)

		return model.Campaing{}, err
	}

	return campaing, nil
}

func (c *Campaing) GetById(ctx context.Context, param uuid.UUID) (model.Campaing, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	tx, err := c.conn.BeginTx(ctx, nil)
	if err != nil {
		easyzap.Warn("select campaing by id %v cancel by context. msg: %v", param, err)

		return model.Campaing{}, err
	}

	var campaing model.Campaing
	row := tx.QueryRowContext(ctx, "select * from campaing where id = $1 and active = true;", param)

	err = row.Scan(&campaing.Id, &campaing.UserId, &campaing.SlugId, &campaing.MerchantId, &campaing.CreatedAt, &campaing.UpdatedAt, &campaing.Active, &campaing.Lat, &campaing.Long, &campaing.Clicks, &campaing.Impressions)
	if err != nil {
		if err == sql.ErrNoRows {
			easyzap.Warn("no lines found for id: %v", param, err)
			tx.Rollback()

			return model.Campaing{}, nil
		}
		easyzap.Errorf("scan campaing %v fail. msg: %v", param, err)
		tx.Rollback()

		return model.Campaing{}, err
	}

	err = tx.Commit()
	if err != nil {
		easyzap.Warn("select campaing by id %v fail. msg: %v", param, err)

		return model.Campaing{}, err
	}

	return campaing, nil
}

func (c *Campaing) Create(ctx context.Context, params model.Campaing) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	tx, err := c.conn.BeginTx(ctx, nil)
	if err != nil {
		easyzap.Warn("create campaing %v cancel by context. msg: %v", params, err)

		return err
	}

	_, err = tx.ExecContext(ctx, "insert into campaing(id,user_id,slug_id,merchant_id,created_at,updated_at,active,lat,long,clicks,impressions) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)",
		params.Id, params.UserId, params.SlugId, params.MerchantId, params.CreatedAt, params.UpdatedAt, params.Active, params.Lat, params.Long, params.Clicks, params.Impressions)
	if err != nil {
		tx.Rollback()
		easyzap.Warn("create campaing %v query exec fail rollbak. msg: %v", params, err)

		return err
	}

	err = tx.Commit()
	if err != nil {
		easyzap.Warn("create campaing %v fail. msg: %v", params, err)

		return err
	}

	return nil
}

//vou usar no producer
/* func (c *Campaing) Select(param uuid.UUID) (model.Campaing, error) {
	span, ctx := tracer.StartSpanFromContext(context.Background(), "campaing_repository.select",
		tracer.ResourceName("postgres"),
		tracer.SpanType("db"),
	)
	defer span.Finish()

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	tx, err := c.conn.BeginTx(ctx, nil)
	if err != nil {
		easyzap.Warn("select campaing id %v cancel by context. msg: %v", param, err)
		span.Finish(tracer.WithError(err))

		return model.Campaing{}, nil
	}

	var campaing model.Campaing
	row, err := tx.QueryContext(ctx, "select * from campaing where id $1", param)
	if err != nil {
		easyzap.Warn("select campaing id %v query exec fail. msg: %v", param, err)
		span.Finish(tracer.WithError(err))

		return model.Campaing{}, nil
	}
	defer row.Close()

	err = row.Scan(&campaing)
	if err != nil {
		tx.Rollback()
		easyzap.Warn("scan campaing id %v fail. msg: %v", param, err)
		span.Finish(tracer.WithError(err))

		return model.Campaing{}, nil
	}

	err = tx.Commit()
	if err != nil {
		easyzap.Warn("select campaing id %v fail. msg: %v", param, err)
		span.Finish(tracer.WithError(err))

		return model.Campaing{}, nil
	}
	return campaing, nil
} */

func (c *Campaing) Update(ctx context.Context, params model.Campaing) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	tx, err := c.conn.BeginTx(ctx, nil)
	if err != nil {
		easyzap.Warn("update campaing %v cancel by context. msg: %v", params, err)

		return err
	}
	result, err := tx.ExecContext(ctx, "update campaing set user_id=$2,slug_id=$3,updated_at=$4,active=$5,lat=$6,long=$7,clicks=$8,impressions=$9 where id = $1",
		params.Id, params.UserId, params.SlugId, params.UpdatedAt, params.Active, params.Lat, params.Long, params.Clicks, params.Impressions)
	if err != nil {
		tx.Rollback()
		easyzap.Warn("update campaing %v query exec fail. msg: %v", params, err)

		return err
	}
	err = tx.Commit()
	if err != nil {
		easyzap.Warn("update campaing %v fail. msg: %v", params, err)

		return err
	}
	result.LastInsertId()

	return nil
}

func (c *Campaing) Delete(ctx context.Context, param uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	tx, err := c.conn.BeginTx(ctx, nil)
	if err != nil {
		easyzap.Warn("delete campaing id %v cancel by context. msg: %v", param, err)

		return err
	}
	_, err = tx.ExecContext(ctx, "delete from campaing where id = $1", param)
	if err != nil {
		tx.Rollback()
		easyzap.Warn("delete campaing id %v fail. msg: %v", param, err)

		return err
	}
	err = tx.Commit()
	if err != nil {
		easyzap.Warn("delete campaing id %v fail. msg: %v", param, err)

		return err
	}
	return nil
}
