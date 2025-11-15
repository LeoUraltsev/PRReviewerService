package pg

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNoConnectionStringENV = errors.New("ENV: POSTGRES_PR_CONNECTION_STRING not found")
)

type Storage struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewStorage(ctx context.Context, log *slog.Logger) (*Storage, error) {
	conn := os.Getenv("POSTGRES_PR_CONNECTION_STRING")
	if conn == "" {
		return nil, ErrNoConnectionStringENV
	}

	cfg, err := pgxpool.ParseConfig(conn)
	if err != nil {
		return nil, fmt.Errorf("failed parse config: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed create pool: %v", err)
	}

	s := &Storage{db: pool, log: log}

	err = s.checkConnection(ctx)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Storage) Close(ctx context.Context) error {
	s.db.Close()
	return nil
}

func (s *Storage) Ping(ctx context.Context) error {
	err := s.db.Ping(ctx)
	if err != nil {
		s.log.Warn("failed ping to postgres", "err", err)
		return fmt.Errorf("failed ping to postgres: %v", err)
	}
	return nil
}

func (s *Storage) checkConnection(ctx context.Context) error {
	s.log.Info("checking connection to postgres")
	err := s.Ping(ctx)
	if err == nil {
		s.log.Info("connected to postgres")
		return nil
	}

	idleDuration := 5 * time.Second
	maxAttemptsRetryConnection := 3
	ticker := time.NewTicker(idleDuration)
	defer ticker.Stop()

	for i := range maxAttemptsRetryConnection {
		s.log.Info(fmt.Sprintf("[%d/%d] trying to connect to postgres", i+1, maxAttemptsRetryConnection))
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			err = s.Ping(ctx)
			if err == nil {
				s.log.Info("connected to postgres")
				return nil
			}

		}

	}
	return err
}
