package pg

import (
	"context"

	"github.com/LeoUraltsev/PRReviewerService/internal/domain"
)

type member struct {
	ID       string
	Username string
	TeamName string
	IsActive bool
}

func (s *Storage) SaveTeam(ctx context.Context, name string) error {

	q := `insert into teams (name) values ($1)`

	_, err := s.db.Exec(ctx, q, name)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetTeam(ctx context.Context, name string) (*domain.Team, error) {
	q := `select name from teams where name = $1`
	row := s.db.QueryRow(ctx, q, name)
	err := row.Scan(&name)
	if err != nil {
		return nil, err
	}
	return &domain.Team{
		TeamName: name,
	}, nil
}
