package domain

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	OPEN   Status = "OPEN"
	CLOSED Status = "CLOSED"
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
