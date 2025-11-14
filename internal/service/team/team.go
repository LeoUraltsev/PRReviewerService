package team

import (
	"context"
	"fmt"

	"github.com/LeoUraltsev/PRReviewerService/internal/domain"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s Service) Save(ctx context.Context, team *domain.Team) error {
	fmt.Println(team)
	return nil
}

func (s Service) Get(ctx context.Context, teamName string) (*domain.Team, error) {
	fmt.Println(teamName)
	return nil, nil
}
