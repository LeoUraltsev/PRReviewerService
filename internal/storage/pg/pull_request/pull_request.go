package pull_request

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/LeoUraltsev/PRReviewerService/internal/domain"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

type pullRequest struct {
	ID                string
	Name              string
	AuthorID          string
	Status            string
	AssignedReviewers []string
	NeedMoreReviewers bool
	CreatedAt         time.Time
	MergedAt          *time.Time
}

func NewStorage(log *slog.Logger, pool *pgxpool.Pool) *Storage {
	return &Storage{
		pool: pool,
		log:  log,
	}
}

func (s *Storage) Save(ctx context.Context, pullRequest *domain.PullRequest) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	q := `
    INSERT INTO pull_requests (
	    id,
        name, 
        author_id, 
        status, 
        need_more_reviewers
    ) VALUES ($1, $2, $3, $4, $5)
    `

	_, err = tx.Exec(ctx, q,
		pullRequest.ID,
		pullRequest.Name,
		pullRequest.AuthorID,
		pullRequest.Status,
		pullRequest.NeedMoreReviewers,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return domain.ErrPRAlreadyExists
		}
		return err
	}

	// Затем добавляем ревьюверов
	if len(pullRequest.AssignedReviewers) > 0 {
		q = `
        INSERT INTO reviewers (pr_id, user_id) 
        VALUES ($1, $2)
        `

		for _, reviewerID := range pullRequest.AssignedReviewers {
			_, err = tx.Exec(context.Background(), q, pullRequest.ID, reviewerID)
			if err != nil {
				return err
			}
		}
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetByID(ctx context.Context, prID string) (*domain.PullRequest, error) {
	q := `SELECT 
    pr.id, 
    pr.name, 
    pr.author_id, 
    pr.status, 
    pr.need_more_reviewers, 
    pr.created_at, 
    pr.merged_at,
    COALESCE(array_agg(DISTINCT rv.user_id) FILTER (WHERE rv.user_id IS NOT NULL), ARRAY[]::text[]) AS reviewers
FROM pull_requests pr
LEFT JOIN reviewers rv ON rv.pr_id = pr.id 
WHERE pr.id = $1
GROUP BY pr.id, pr.name, pr.author_id, pr.status, pr.need_more_reviewers, pr.created_at, pr.merged_at
`
	var pr pullRequest
	err := s.pool.QueryRow(ctx, q, prID).Scan(
		&pr.ID,
		&pr.Name,
		&pr.AuthorID,
		&pr.Status,
		&pr.NeedMoreReviewers,
		&pr.CreatedAt,
		&pr.MergedAt,
		&pr.AssignedReviewers,
	)
	if err != nil {
		return nil, err
	}

	return toDomainPullRequest(&pr), nil
}

func (s *Storage) UpdateStatus(ctx context.Context, id string, status domain.Status) (*domain.PullRequest, error) {
	q := ""
	if status == domain.Merged {
		q = `UPDATE pull_requests SET status = $1, merged_at = timezone('utc', now()) WHERE id = $2`
	} else {
		q = `UPDATE pull_requests SET status = $1 WHERE id = $2`
	}
	_, err := s.pool.Exec(ctx, q, status, id)
	if err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id)
}

func (s *Storage) Reassign(ctx context.Context, prID string, oldUserID string, newUserID string) (*domain.PullRequest, error) {
	q := `update reviewers set user_id = $1 where pr_id = $2 and user_id = $3`
	_, err := s.pool.Exec(ctx, q, newUserID, prID, oldUserID)
	if err != nil {
		return nil, err
	}
	return s.GetByID(ctx, prID)
}

func (s *Storage) GetPRByUserID(ctx context.Context, userID string) ([]*domain.PullRequest, error) {
	q := `SELECT pr_id FROM reviewers WHERE user_id = $1`
	rows, err := s.pool.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var prs []*domain.PullRequest
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		pr, err := s.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}
	return prs, nil
}

func toDomainPullRequest(pr *pullRequest) *domain.PullRequest {

	return &domain.PullRequest{
		ID:                pr.ID,
		Name:              pr.Name,
		AuthorID:          pr.AuthorID,
		Status:            domain.Status(pr.Status),
		AssignedReviewers: pr.AssignedReviewers,
		NeedMoreReviewers: pr.NeedMoreReviewers,
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}
}
