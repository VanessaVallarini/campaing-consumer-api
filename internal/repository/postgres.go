package repository

import (
	"campaing-comsumer-service/internal/config"
	"database/sql"

	"github.com/lockp111/go-easyzap"

	_ "github.com/lib/pq"
)

func NewPostgresClient(cfg config.DatabaseConfig) *sql.DB {
	conn, err := sql.Open(cfg.PostgresDriver, cfg.DatabaseConnStr)
	if err != nil {
		easyzap.Fatal("new db client error: %v", err)
	}

	return conn
}
