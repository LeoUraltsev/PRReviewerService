package pull_request

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/LeoUraltsev/PRReviewerService/internal/domain"
)

/*
Возможные проблемы
1. Пользователя которого мы достали для замены мог взять другой pr, нужна транзакция для изоляции
2. Нейминг интерфейсов пошел погулять
*/

type RepoPR interface {
	Save(ctx context.Context, pullRequest *domain.PullRequest) error
	GetByID(ctx context.Context, prID string) (*domain.PullRequest, error)
	UpdateStatus(ctx context.Context, id string, status domain.Status) (*domain.PullRequest, error)
	Reassign(ctx context.Context, prID string, oldUserID string, newUserID string) (*domain.PullRequest, error)
}

type UserRepo interface {
	GetInactiveUsers(ctx context.Context, teamName string, limit int, excludeUsers []string) ([]*domain.User, error)
	GetByID(ctx context.Context, userID string) (*domain.User, error)
	UpdateIsActive(ctx context.Context, userID string, active bool) (*domain.User, error)
}

type Service struct {
	repoPR   RepoPR
	repoUser UserRepo
}

func NewService(pr RepoPR, user UserRepo) *Service {
	return &Service{
		repoPR:   pr,
		repoUser: user,
	}
}

func (s *Service) SavePullRequest(ctx context.Context, prID string, prName string, authorID string) (*domain.PullRequest, error) {
	user, err := s.repoUser.GetByID(ctx, authorID)
	if err != nil {
		return nil, err
	}

	slog.Info("call SavePullRequest")
	needMoreReviewers := false
	countUsers := 2
	excludeUsers := []string{authorID}
	users, err := s.repoUser.GetInactiveUsers(ctx, user.TeamName, countUsers, excludeUsers)
	if err != nil {
		return nil, err
	}

	if len(users) < countUsers {
		needMoreReviewers = true
	}

	reviewersIDs := make([]string, 0, len(users))
	for _, v := range users {
		_, err = s.repoUser.UpdateIsActive(ctx, v.UserID, false)
		if err != nil {
			slog.Warn("failed to update isActive for PR reviewer", "user", v.UserID)
		}
		reviewersIDs = append(reviewersIDs, v.UserID)
	}

	err = s.repoPR.Save(ctx, &domain.PullRequest{
		ID:                prID,
		Name:              prName,
		AuthorID:          authorID,
		Status:            domain.Open,
		AssignedReviewers: reviewersIDs,
		NeedMoreReviewers: needMoreReviewers,
	})
	if err != nil {
		return nil, err
	}
	slog.Info("save pull request")
	slog.Info("get pr by id")
	pr, err := s.repoPR.GetByID(ctx, prID)
	if err != nil {
		slog.Error("get pr by id failed", err)
		return nil, err
	}
	return pr, nil
}

func (s *Service) MergePullRequest(ctx context.Context, prID string) (*domain.PullRequest, error) {
	pr, err := s.repoPR.GetByID(ctx, prID)
	if err != nil {
		return nil, err
	}

	if pr.Status == domain.Merged {
		return pr, nil
	}

	pr.ChangeStatus(domain.Merged)
	pr, err = s.repoPR.UpdateStatus(ctx, pr.ID, pr.Status)
	if err != nil {
		return nil, err
	}
	return pr, nil
}

func (s *Service) ReassignReviewerPullRequest(ctx context.Context, prID string, reviewerID string) (*domain.PullRequest, error) {
	pr, err := s.repoPR.GetByID(ctx, prID)
	if err != nil {
		return nil, err
	}

	if pr.Status == domain.Merged {
		return nil, domain.ErrReassignPRMerged
	}

	found := false
	for _, assignedID := range pr.AssignedReviewers {
		if assignedID == reviewerID {
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("reviewer %s not found in PR %s", reviewerID, pr.ID)
	}

	user, err := s.repoUser.GetByID(ctx, reviewerID)
	if err != nil {
		return nil, err
	}
	excludeUsers := []string{pr.AuthorID, user.UserID}
	newUser, err := s.repoUser.GetInactiveUsers(ctx, user.TeamName, 1, excludeUsers)
	if err != nil {
		return nil, err
	}
	if len(newUser) < 1 {
		return nil, fmt.Errorf("no inactive user")
	}
	u := newUser[0]

	user.ChangeActive(true)
	u.ChangeActive(false)

	newPR, err := s.repoPR.Reassign(ctx, prID, user.UserID, u.UserID)
	if err != nil {
		return nil, err
	}
	_, err = s.repoUser.UpdateIsActive(ctx, user.UserID, user.IsActive)
	if err != nil {
		return nil, err
	}
	_, err = s.repoUser.UpdateIsActive(ctx, u.UserID, u.IsActive)
	if err != nil {
		return nil, err
	}

	return newPR, nil
}
