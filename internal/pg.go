package internal

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresPool struct {
	pool *pgxpool.Pool
}

func NewPostgres(connString string) (*PostgresPool, error) {
	pgxPool, err := pgxpool.New(context.TODO(), connString)
	if err != nil {
		return nil, err
	}
	return &PostgresPool{pool: pgxPool}, nil
}

func (p *PostgresPool) ReadHelloWorld() (greeting string, err error) {
	err = p.pool.QueryRow(context.TODO(), "select 'Hello, world!'").Scan(&greeting)
	return
}

func (p *PostgresPool) Close() {
	p.pool.Close()
}
