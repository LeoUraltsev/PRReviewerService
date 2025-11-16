package pg

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/LeoUraltsev/PRReviewerService/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	Pool *pgxpool.Pool
	log  *slog.Logger
	cfg  *config.Config
}

func NewStorage(ctx context.Context, log *slog.Logger, cfg *config.Config) (*Storage, error) {
	conn := cfg.ConnectionString

	pgCfg, err := pgxpool.ParseConfig(conn)
	if err != nil {
		return nil, fmt.Errorf("failed parse config: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgCfg)
	if err != nil {
		return nil, fmt.Errorf("failed create pool: %v", err)
	}

	s := &Storage{
		Pool: pool,
		log:  log,
		cfg:  cfg,
	}

	err = s.checkConnection(ctx)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Storage) Close(ctx context.Context) error {
	s.Pool.Close()
	return nil
}

func (s *Storage) Ping(ctx context.Context) error {
	err := s.Pool.Ping(ctx)
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

	idleDuration := s.cfg.RetryInterval
	maxAttemptsRetryConnection := s.cfg.MaxRetries
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
