package infra

import (
    "context"
    "fmt"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/Rainminds/gantral/infra/db"
)

type Store struct {
    *db.Queries
    pool *pgxpool.Pool
}

func NewStore(ctx context.Context, dsn string) (*Store, error) {
    pool, err := pgxpool.New(ctx, dsn)
    if err != nil {
        return nil, fmt.Errorf("unable to create connection pool: %w", err)
    }

    if err := pool.Ping(ctx); err != nil {
        return nil, fmt.Errorf("unable to ping database: %w", err)
    }

    return &Store{
        Queries: db.New(pool),
        pool:    pool,
    }, nil
}

func (s *Store) Close() {
    s.pool.Close()
}
