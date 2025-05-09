package database

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type DB struct {
	Conn *sqlx.DB
}

func (d *DB) Close() error {
	return d.Conn.Close()
}

func (d *DB) BeginTxx(ctx context.Context) (TxInterface, error) {
	tx, err := d.Conn.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &Tx{tx: tx}, nil
}


func (d *DB) RunWithTx(ctx context.Context, fn func(tx TxInterface) error) error {
	tx, err := d.BeginTxx(ctx)
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		_ = tx.Rollback
		return err
	}

	return tx.Commit()
}
	