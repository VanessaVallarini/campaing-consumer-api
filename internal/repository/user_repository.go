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

type User struct {
	conn    *sql.DB
	metrics *metrics.Metrics
}

func NewUserRepository(metrics *metrics.Metrics, conn *sql.DB) *User {
	return &User{
		conn:    conn,
		metrics: metrics,
	}
}

func (c *User) GetById(ctx context.Context, param uuid.UUID) (model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	tx, err := c.conn.BeginTx(ctx, nil)
	if err != nil {
		mv := []string{"GetById", "error", "starts_transaction"}
		c.metrics.UserRepository.WithLabelValues(mv...).Inc()
		errWrap := errors.Wrapf(err, "select user %v cancel by context", param)
		easyzap.Warn(ctx, errWrap, "select user cancel by context")
		return model.User{}, errWrap
	}

	var user model.User
	query := "select * from " + "\"user\"" + " where id = $1"
	row := tx.QueryRowContext(ctx, query, param)

	err = row.Scan(&user.Id, &user.Email, &user.CreatedAt, &user.UpdatedAt, &user.Active)
	if err != nil {
		if err == sql.ErrNoRows {
			tx.Rollback()
			return model.User{}, nil
		}
		tx.Rollback()
		mv := []string{"GetById", "error", "scan"}
		c.metrics.UserRepository.WithLabelValues(mv...).Inc()
		errWrap := errors.Wrapf(err, "scan user %v fail", param)
		easyzap.Error(ctx, errWrap, "scan user fail")
		return model.User{}, errWrap
	}

	err = tx.Commit()
	if err != nil {
		mv := []string{"GetById", "error", "commit"}
		c.metrics.UserRepository.WithLabelValues(mv...).Inc()
		errWrap := errors.Wrapf(err, "select user id %v fail", param)
		easyzap.Error(ctx, errWrap, "select user fail")
		return model.User{}, errWrap
	}

	mv := []string{"GetById", "success", ""}
	c.metrics.UserRepository.WithLabelValues(mv...).Inc()

	return user, nil
}
