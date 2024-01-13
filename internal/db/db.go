package db

import (
	"campaing-comsumer-service/internal/config"
	"context"
	"database/sql"

	"github.com/joomcode/errorx"
	_ "github.com/lib/pq"
	"github.com/lockp111/go-easyzap"
)

type IDb interface {
	Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
	Ping(ctx context.Context) *errorx.Error
	Close()
}

type Db struct {
	conn *sql.DB
}

func NewDb(ctx context.Context, cfg config.Config) (*Db, error) {
	db, err := sql.Open(cfg.DatabaseConfig.PostgresDriver, cfg.DatabaseConfig.DatabaseConnStr)
	if err != nil {
		return nil, err
	}

	return &Db{
		conn: db,
	}, nil
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
