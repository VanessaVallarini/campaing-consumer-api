package db

import (
	"campaing-comsumer-service/internal/config"
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/lockp111/go-easyzap"
)

type IDb interface {
	Ping(ctx context.Context) error
	Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
	Close()
}

type Db struct {
	conn *sql.DB
}

func NewDb(cfg config.Config) *Db {
	db, err := sql.Open(cfg.DatabaseConfig.PostgresDriver, cfg.DatabaseConfig.DatabaseConnStr)
	if err != nil {
		easyzap.Panicf("configuration error: %v", err)
	}

	return &Db{
		conn: db,
	}
}

func (db *Db) Ping(ctx context.Context) error {
	err := db.conn.Ping()
	if err != nil {
		easyzap.Error(ctx, err, "failed to ping database")

		return err
	}

	return nil
}

func (db *Db) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return db.conn.Exec(query, args)
}

func (db *Db) Close() {
	db.conn.Close()
}
