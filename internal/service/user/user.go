package user

import (
	"context"
	"fmt"
	"time"

	"github.com/LeoUraltsev/PRReviewerService/internal/domain"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s Service) GetUserPullRequest(ctx context.Context, userId string) ([]*domain.PullRequest, error) {
	fmt.Println(userId)
	return []*domain.PullRequest{
		{
			ID:                "1",
			Name:              "Name 1",
			AuthorID:          "asdasd",
			Status:            domain.Open,
			AssignedReviewers: nil,
			NeedMoreReviewers: false,
			CreatedAt:         time.Time{},
			MergedAt:          time.Time{},
		},
	}, nil
}

func (s Service) UpdateIsActive(ctx context.Context, userId string, isActive bool) (*domain.User, error) {
	fmt.Println(userId, isActive)
	return &domain.User{
		UserID:   "1",
		Username: "Name 1",
		TeamName: "Team",
		IsActive: false,
	}, nil
}
