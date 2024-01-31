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

type User struct {
	conn *sql.DB
}

func NewUserRepository(conn *sql.DB) *User {
	return &User{
		conn: conn,
	}
}

func (c *User) GetById(param uuid.UUID) (model.User, error) {
	span, ctx := tracer.StartSpanFromContext(context.Background(), "user_repository.get-by-id",
		tracer.ResourceName("postgres"),
		tracer.SpanType("db"),
	)
	defer span.Finish()

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	tx, err := c.conn.BeginTx(ctx, nil)
	if err != nil {
		easyzap.Warn("select user %v cancel by context. msg: %v", param, err)
		span.Finish(tracer.WithError(err))

		return model.User{}, err
	}

	var user model.User
	query := "select * from " + "\"user\"" + " where id = $1"
	row := tx.QueryRowContext(ctx, query, param)

	err = row.Scan(&user.Id, &user.Email, &user.CreatedAt, &user.UpdatedAt, &user.Active)
	if err != nil {
		if err == sql.ErrNoRows {
			easyzap.Warn("no lines found for user_id: %v", param, err)
			tx.Rollback()

			return model.User{}, nil

		}
		easyzap.Errorf("scan user %v fail. msg: %v", param, err)
		tx.Rollback()
		span.Finish(tracer.WithError(err))

		return model.User{}, err
	}

	err = tx.Commit()
	if err != nil {
		easyzap.Warn("select user %v fail. msg: %v", param, err)
		span.Finish(tracer.WithError(err))

		return model.User{}, err
	}

	return user, nil
}
