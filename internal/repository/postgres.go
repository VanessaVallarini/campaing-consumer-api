package repository

import (
	"campaing-comsumer-service/internal/config"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

func CreatePool(ctx context.Context, config *config.DatabaseConfig) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(fmt.Sprintf("host=%s port=%d user=%s dbname=%s "+
		"password=%s sslmode=%s pool_min_conns=%d pool_max_conns=%d "+
		"pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s default_query_exec_mode=cache_describe",
		config.Host,
		config.Port,
		config.Username,
		config.Database,
		config.Password,
		"disable",
		config.Conn.Min,
		config.Conn.Max,
		config.Conn.Lifetime,
		config.Conn.IdleTime,
	))

	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(
		ctx,
		cfg,
	)
	return pool, err
}
