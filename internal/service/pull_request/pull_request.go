package pull_request

import (
	"context"

	"github.com/LeoUraltsev/PRReviewerService/internal/domain"
)

/*
Возможные проблемы
1. Пользователя которого мы достали для замены мог взять другой pr, нужна транзакция для изоляции
2. Нейминг интерфейсов пошел погулять
*/

type RepoPR interface {
	Save(ctx context.Context, pullRequest *domain.PullRequest) (*domain.PullRequest, error)
	GetByID(ctx context.Context, prID string) (*domain.PullRequest, error)
	Update(ctx context.Context, pr *domain.PullRequest) (*domain.PullRequest, error)
	Reassign(ctx context.Context, pr *domain.PullRequest) (*domain.PullRequest, error)
}

type TeamRepo interface {
	GetFirstInactiveUser(ctx context.Context) (*domain.User, error)
}

type UserRepo interface {
	GetByID(ctx context.Context, userID string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
}

type Service struct {
	repoPR   RepoPR
	repoTeam TeamRepo
	repoUser UserRepo
}

func NewService(pr RepoPR, team TeamRepo, user UserRepo) *Service {
	return &Service{
		repoPR:   pr,
		repoTeam: team,
		repoUser: user,
	}
}

func (s *Service) SavePullRequest(ctx context.Context, prID string, prName string, authorID string) (*domain.PullRequest, error) {
	pr, err := s.repoPR.Save(ctx, &domain.PullRequest{
		ID:       prID,
		Name:     prName,
		AuthorID: authorID,
	})
	if err != nil {
		return nil, err
	}
	return pr, nil
}

func (s *Service) MergePullRequest(ctx context.Context, prID string) (*domain.PullRequest, error) {
	pr, err := s.repoPR.GetByID(ctx, prID)
	if err != nil {
		return nil, err
	}

	if pr.Status == domain.Closed {
		return pr, nil
	}

	pr.ChangeStatus(domain.Open)

	pr, err = s.repoPR.Update(ctx, pr)
	if err != nil {
		return nil, err
	}
	return pr, nil
}

//todo: обязательно в транзакцию

func (s *Service) ReassignReviewerPullRequest(ctx context.Context, prID string, reviewerID string) (*domain.PullRequest, error) {
	pr, err := s.repoPR.GetByID(ctx, prID)
	if err != nil {
		return nil, err
	}

	if pr.Status == domain.Closed {
		return pr, domain.ErrReassignPRMerged
	}

	user, err := s.repoUser.GetByID(ctx, reviewerID)

	newUser, err := s.repoTeam.GetFirstInactiveUser(ctx)
	if err != nil {
		return nil, err
	}

	newUser.ChangeActive(false)
	err = s.repoUser.Update(ctx, user)

	for _, r := range pr.AssignedReviewers {
		if r == reviewerID {
			r = newUser.UserID
		}
	}

	newUser.ChangeActive(true)
	err = s.repoUser.Update(ctx, newUser)
	if err != nil {
		return nil, err
	}
	return pr, nil
}
