package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
	"url-storage/internal/config"
	"url-storage/internal/storage"
)

type UrlStorage struct {
	pool    *pgxpool.Pool
	timeout time.Duration
}

func New(cfg *config.PostgresConfig) (*UrlStorage, error) {
	const op = "postgres.New"

	connString := fmt.Sprintf(
		"user=%s password=%s host=%s dbname=%s sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.DBName,
		cfg.SSLMode,
	)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	createTableStmt := `
		CREATE TABLE IF NOT EXISTS urls (
			id SERIAL PRIMARY KEY,
			url TEXT NOT NULL,
			alias TEXT NOT NULL
		);`

	createIndexStmt := `
		CREATE UNIQUE INDEX IF NOT EXISTS alias_idx ON urls (alias);`

	batch := &pgx.Batch{}
	batch.Queue(createTableStmt)
	batch.Queue(createIndexStmt)

	ctx, cancel = context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	_, err = pool.SendBatch(ctx, batch).Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &UrlStorage{
		pool:    pool,
		timeout: cfg.Timeout,
	}, nil
}

func (s *UrlStorage) Insert(ctx context.Context, url, alias string) (int64, error) {
	const op = "postgres.Insert"

	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var id int64
	err := s.pool.QueryRow(ctx, "INSERT INTO urls (url, alias) VALUES ($1, $2) RETURNING id;", url, alias).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		ok := errors.As(err, &pgErr)
		if ok {
			if pgErr.Code == "23505" { // duplicate key value violates unique constraint
				return -1, storage.ErrAlreadyExist
			}
		}
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *UrlStorage) Get(ctx context.Context, alias string) (string, error) {
	const op = "postgres.Get"

	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var url string
	err := s.pool.QueryRow(ctx, "SELECT url FROM urls WHERE alias = $1;", alias).Scan(&url)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return url, nil
}
