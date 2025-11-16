package user

import (
	"context"

	"github.com/LeoUraltsev/PRReviewerService/internal/domain"
)

//todo: логирование

type RepoPR interface {
	GetPRByUserID(ctx context.Context, userID string) ([]*domain.PullRequest, error)
}

type RepoUsers interface {
	CheckExists(ctx context.Context, userID string) error
	UpdateIsActive(ctx context.Context, userID string, active bool) (*domain.User, error)
}

type Service struct {
	repoUsers RepoUsers
	repoPR    RepoPR
}

func NewService(pr RepoPR, users RepoUsers) *Service {
	return &Service{
		repoPR:    pr,
		repoUsers: users,
	}
}

func (s *Service) GetUserPullRequest(ctx context.Context, userId string) ([]*domain.PullRequest, error) {
	err := s.repoUsers.CheckExists(ctx, userId)
	if err != nil {
		return nil, err
	}

	pr, err := s.repoPR.GetPRByUserID(ctx, userId)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *Service) UpdateIsActive(ctx context.Context, userId string, isActive bool) (*domain.User, error) {
	user, err := s.repoUsers.UpdateIsActive(ctx, userId, isActive)
	if err != nil {
		return nil, err
	}
	return user, nil
}
