package db

import (
	"campaing-comsumer-service/internal/config"
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx"
	"github.com/joomcode/errorx"
	_ "github.com/lib/pq"
	"github.com/lockp111/go-easyzap"
)

type IDb interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Ping(ctx context.Context) *errorx.Error
	Close()
	Read()
}

type Db struct {
	conn *sql.DB
}

func NewDb(ctx context.Context, cfg config.Config) (*Db, error) {
	db, err := sql.Open(cfg.Database.PostgresDriver, cfg.Database.DatabaseConnStr)
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

func (db *Db) Begin() (*sql.Tx, error) {
	return db.conn.Begin()
}

func (db *Db) Read() {
	_, err := db.conn.Exec("insert into campaing(id,user_id,slug_id,merchant_id,active,lat,long) "+
		"VALUES($1,$2,$3,$4,$5,$6,$7)", "cd1d5c2f-716f-42fb-a6d6-27031de0ae67", "cd1d5c2f-716f-42fb-a6d6-27031de0ae67", "cd1d5c2f-716f-42fb-a6d6-27031de0ae67", "cd1d5c2f-716f-42fb-a6d6-27031de0ae67", true, 45.6085, -73.5493)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("value inserted")
	}
}

func (db *Db) Close() {
	db.conn.Close()
}
