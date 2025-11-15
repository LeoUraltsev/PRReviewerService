package team

import (
	"context"

	"github.com/LeoUraltsev/PRReviewerService/internal/domain"
)

type RepoTeam interface {
	Save(ctx context.Context, teamName string) error
	CheckTeam(ctx context.Context, teamName string) error
}

type RepoUser interface {
	SaveBatch(ctx context.Context, user []*domain.User) error
	GetUsersByTeamName(ctx context.Context, teamName string) ([]*domain.User, error)
}

type Service struct {
	repo     RepoTeam
	repoUser RepoUser
}

func NewService(user RepoUser, team RepoTeam) *Service {
	return &Service{
		repo:     team,
		repoUser: user,
	}
}

// todo: объеденить в транзакцию
func (s Service) Save(ctx context.Context, team *domain.Team) error {
	err := s.repo.Save(ctx, team.TeamName)
	if err != nil {
		return err
	}

	err = s.repoUser.SaveBatch(ctx, team.Members)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) Get(ctx context.Context, teamName string) (*domain.Team, error) {
	err := s.repo.CheckTeam(ctx, teamName)
	if err != nil {
		return nil, err
	}

	members, err := s.repoUser.GetUsersByTeamName(ctx, teamName)
	if err != nil {
		return nil, err
	}

	team := &domain.Team{
		TeamName: teamName,
		Members:  members,
	}

	return team, nil
}
