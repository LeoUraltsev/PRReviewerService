package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/LeoUraltsev/PRReviewerService/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Storage) CreateUser(ctx context.Context, user *domain.User) error {
	q := `INSERT INTO users (id, username, team_name, is_active) values ($1, $2, $3, $4)`

	_, err := s.db.Exec(ctx, q, user.UserID, user.Username, user.TeamName, user.IsActive)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf("user %s already exists", user.UserID)
			}
		}
		return err
	}
	return nil
}

func (s *Storage) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	q := `SELECT id, username, team_name, is_active FROM users WHERE id = $1`
	row := s.db.QueryRow(ctx, q, userID)
	user := &domain.User{}
	err := row.Scan(user)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *Storage) UpdateUserActive(ctx context.Context, id string, isActive bool) error {
	q := `UPDATE users SET is_active = $1 WHERE id = $2`
	_, err := s.db.Exec(ctx, q, isActive, id)
	if err != nil {
		return err
	}
	return nil
}
