package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

type Postgres struct {
	Pool *pgxpool.Pool
}

func New(url string) (*Postgres, error) {
	pg := &Postgres{}

	// Логика Retry: пытаемся подключиться несколько раз, пока Docker-контейнер БД просыпается
	for i := 0; i < _defaultConnAttempts; i++ {
		pool, err := pgxpool.New(context.Background(), url)
		if err == nil {
			pg.Pool = pool
			return pg, nil
		}

		fmt.Printf("Postgres is trying to connect, attempt %d...\n", i+1)
		time.Sleep(_defaultConnTimeout)
	}

	return nil, fmt.Errorf("postgres - New - conn attempts exceeded")
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
