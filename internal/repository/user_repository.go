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
		easyzap.Warn("select user %v cancel by context. msg: %v", param, err)
		mv := []string{"GetById", "error", "starts_transaction"}
		c.metrics.UserRepository.WithLabelValues(mv...).Inc()

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
			mv := []string{"GetById", "error", "no_rows"}
			c.metrics.UserRepository.WithLabelValues(mv...).Inc()

			return model.User{}, nil

		}
		easyzap.Errorf("scan user %v fail. msg: %v", param, err)
		tx.Rollback()
		mv := []string{"GetById", "error", "scan"}
		c.metrics.UserRepository.WithLabelValues(mv...).Inc()

		return model.User{}, err
	}

	err = tx.Commit()
	if err != nil {
		easyzap.Warn("select user %v fail. msg: %v", param, err)
		mv := []string{"GetById", "error", "commit"}
		c.metrics.UserRepository.WithLabelValues(mv...).Inc()

		return model.User{}, err
	}

	mv := []string{"GetById", "success", ""}
	c.metrics.UserRepository.WithLabelValues(mv...).Inc()

	return user, nil
}
