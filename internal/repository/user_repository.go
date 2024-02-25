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

type User struct {
	db      *pgxpool.Pool
	metrics *metrics.Metrics
}

func NewUserRepository(metrics *metrics.Metrics, db *pgxpool.Pool) *User {
	return &User{
		db:      db,
		metrics: metrics,
	}
}

func (c *User) GetById(ctx context.Context, param uuid.UUID) (model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	var user model.User
	query := "select * from " + "\"user\"" + " where id = $1"
	rows := c.db.QueryRow(ctx, query, param)

	err := rows.Scan(&user.Id, &user.Email, &user.CreatedAt, &user.UpdatedAt, &user.Active)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return model.User{}, nil
		}
		mv := []string{"GetById", "error", "scan"}
		c.metrics.UserRepository.WithLabelValues(mv...).Inc()
		errWrap := errors.Wrapf(err, "scan user %v fail", param)
		easyzap.Errorf("scan user %v fail. msg: %v", param, errWrap)
		return model.User{}, errWrap
	}

	mv := []string{"GetById", "success", ""}
	c.metrics.UserRepository.WithLabelValues(mv...).Inc()

	return user, nil
}
