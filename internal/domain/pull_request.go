package domain

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	Open   Status = "OPEN"
	Closed Status = "CLOSED"
)

type PullRequest struct {
	ID                uuid.UUID
	Name              string
	AuthorID          uuid.UUID
	Status            Status
	AssignedReviewers []uuid.UUID
	NeedMoreReviewers bool
	CreatedAt         time.Time
	MergedAt          time.Time
}
