package user

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/LeoUraltsev/PRReviewerService/internal/domain"
	"github.com/jackc/pgx/v5"
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

type user struct {
	ID        string
	Username  string
	TeamName  string
	IsActive  bool
	CreatedAt time.Time
}

func (s *Storage) GetByID(ctx context.Context, userID string) (*domain.User, error) {
	q := `select id, username, team_name, is_active, created_at from users where id = $1`
	var u user
	err := s.pool.QueryRow(ctx, q, userID).Scan(&u.ID, &u.Username, &u.TeamName, &u.IsActive, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return toDomainUser(&u), nil

}

func (s *Storage) UpdateIsActive(ctx context.Context, userID string, active bool) (*domain.User, error) {
	q := `update users set is_active = $1 where id = $2 RETURNING id, username, team_name, is_active, created_at`
	var u user
	err := s.pool.QueryRow(ctx, q, active, userID).Scan(&u.ID, &u.Username, &u.TeamName, &u.IsActive, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return toDomainUser(&u), nil
}

func (s *Storage) SaveUsers(ctx context.Context, users []*domain.User) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	q := `insert into users (id, username, team_name, is_active) values ($1, $2, $3, $4)`

	for _, u := range users {
		_, err = tx.Exec(ctx, q, u.UserID, u.Username, u.TeamName, u.IsActive)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetUsersByTeamName(ctx context.Context, teamName string) ([]*domain.User, error) {
	q := `select id, username, team_name, is_active, created_at from users where team_name = $1`
	rows, err := s.pool.Query(ctx, q, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]*domain.User, 0)
	for rows.Next() {
		var u user
		err = rows.Scan(&u.ID, &u.Username, &u.TeamName, &u.IsActive, &u.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, toDomainUser(&u))
	}
	return users, nil
}

func (s *Storage) CheckExists(ctx context.Context, userID string) error {
	q := `select count(*) from users where id = $1`
	var count int
	err := s.pool.QueryRow(ctx, q, userID).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (s *Storage) GetInactiveUsers(ctx context.Context, teamName string, limit int, excludeUsers []string) ([]*domain.User, error) {
	q := `SELECT id, username, team_name, is_active, created_at from users where team_name = $1 and is_active = true and id != ALL($3) limit $2`

	rows, err := s.pool.Query(ctx, q, teamName, limit, excludeUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]*domain.User, 0)
	for rows.Next() {
		var u user
		err = rows.Scan(&u.ID, &u.Username, &u.TeamName, &u.IsActive, &u.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, toDomainUser(&u))
	}
	return users, nil
}

func toDomainUser(u *user) *domain.User {
	return &domain.User{
		UserID:   u.ID,
		Username: u.Username,
		TeamName: u.TeamName,
		IsActive: u.IsActive,
	}
}
