package posgtres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"url-shortener/internal/config"
	"url-shortener/internal/storage"

	//"database/sql"
	"github.com/jackc/pgconn"
)

type Storage struct {
	db *pgxpool.Pool
}

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func New(ctx context.Context, cs config.Settings) (*Storage, error) {
	var pool *pgxpool.Pool
	var err error
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cs.Username, cs.Password, cs.Host, cs.Port, cs.Database)
	//err = utils.DoWithTries(func() error {
	//	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	//	defer cancel()

	pool, err = pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}
	_, err = pool.Exec(ctx, `
CREATE TABLE IF NOT EXISTS url (
    id SERIAL PRIMARY KEY,
    alias TEXT NOT NULL UNIQUE,
    url TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
`)
	if err != nil {
		return nil, err
	}
	return &Storage{db: pool}, err
}

func (s *Storage) SaveURL(urlToSave, alias string) (int64, error) {
	var id int64
	err := s.db.QueryRow(context.Background(), "INSERT INTO url(url, alias) VALUES ($1, $2) RETURNING id", urlToSave, alias).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, storage.ErrURLExists
		}
		return 0, err
	}
	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	var resURL string

	err := s.db.QueryRow(context.Background(), "SELECT url FROM url WHERE alias = $1", alias).Scan(&resURL)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}
		return "", err
	}
	return resURL, nil
}

func (s *Storage) DeleteURL(alias string) error {
	_, err := s.db.Exec(context.Background(), "DELETE FROM url WHERE alias = $1", alias)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.ErrURLNotFound
		}
		return err
	}
	return nil
}
