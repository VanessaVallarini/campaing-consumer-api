package repository

import (
	"campaing-comsumer-service/internal/metrics"
	"context"

	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lockp111/go-easyzap"
	"github.com/pkg/errors"
)

type Transaction struct {
	db      *pgxpool.Pool
	metrics *metrics.Metrics
}

func NewTransactionDao(metrics *metrics.Metrics, db *pgxpool.Pool) *Transaction {
	return &Transaction{
		db:      db,
		metrics: metrics,
	}
}

func (t *Transaction) Begin(ctx context.Context) (pgx.Tx, error) {
	b, err := t.db.Begin(ctx)
	if err != nil {
		mv := []string{"Begin", "error", ""}
		t.metrics.Transaction.WithLabelValues(mv...).Inc()
		errWrap := errors.Wrap(err, "begin transaction fail")
		easyzap.Error("begin transaction fail")
		return b, errWrap
	}
	return b, err
}

func (t *Transaction) Rollback(ctx context.Context, tx pgx.Tx) error {
	err := tx.Rollback(ctx)
	if err != nil {
		mv := []string{"Rollback", "error", ""}
		t.metrics.Transaction.WithLabelValues(mv...).Inc()
		errWrap := errors.Wrap(err, "rollback transaction fail")
		easyzap.Error("rollback transaction fail")
		return errWrap
	}
	return err
}

func (t *Transaction) Commit(ctx context.Context, tx pgx.Tx) error {
	err := tx.Commit(ctx)
	if err != nil {
		mv := []string{"Commit", "error", ""}
		t.metrics.Transaction.WithLabelValues(mv...).Inc()
		errWrap := errors.Wrap(err, "commit transaction fail")
		easyzap.Error("commit transaction fail")
		return errWrap
	}
	return err
}
