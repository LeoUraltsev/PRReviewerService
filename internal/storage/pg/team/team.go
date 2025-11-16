package team

import (
	"context"
	"errors"
	"log/slog"

	"github.com/LeoUraltsev/PRReviewerService/internal/domain"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

func NewStorage(log *slog.Logger, pool *pgxpool.Pool) *Storage {
	return &Storage{
		log:  log,
		pool: pool,
	}
}

func (s *Storage) Save(ctx context.Context, teamName string) error {
	q := `INSERT INTO teams (name) VALUES ($1)`
	_, err := s.pool.Exec(ctx, q, teamName)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return domain.ErrTeamExists
		}
		return err
	}
	return nil
}

func (s *Storage) CheckExistsTeam(ctx context.Context, teamName string) (bool, error) {
	q := `SELECT COUNT(*) FROM teams WHERE name = $1`
	var count int
	if err := s.pool.QueryRow(ctx, q, teamName).Scan(&count); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return count > 0, nil
}
