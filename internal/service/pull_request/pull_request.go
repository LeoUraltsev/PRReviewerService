package pull_request

import (
	"context"

	"github.com/LeoUraltsev/PRReviewerService/internal/domain"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s Service) SavePullRequest(ctx context.Context, prID string, prName string, authorID string) (*domain.PullRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) MergePullRequest(ctx context.Context, prID string) (*domain.PullRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) ReassignReviewerPullRequest(ctx context.Context, prID string, reviewerID string) (*domain.PullRequest, error) {
	//TODO implement me
	panic("implement me")
}
