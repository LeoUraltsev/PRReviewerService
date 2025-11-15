package pg

import (
	"context"
	"log/slog"

	"github.com/LeoUraltsev/PRReviewerService/internal/domain"
)

func (s *Storage) SaveTeam(ctx context.Context) error {

	return nil
}

func (s *Storage) GetTeamByNameWithMember(ctx context.Context, name string) (*domain.Team, error) {
	s.log.Info("getting team by name", slog.String("name", name))

	q := `SELECT user_id, user_name, team_name, is_active FROM team JOIN users ON team.name = users.team_name WHERE team.name = $1;`

	rows, err := s.db.Query(ctx, q, name)
	if err != nil {
		s.log.Info("error getting team by name", slog.String("name", name), slog.String("err", err.Error()))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan()
	}

	return &domain.Team{
		TeamName: "",
	}, nil
}
