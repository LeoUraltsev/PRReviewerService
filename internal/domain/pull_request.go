package domain

import (
	"time"
)

type Status string

const (
	Open   Status = "OPEN"
	Closed Status = "CLOSED"
)

type PullRequest struct {
	ID                string
	Name              string
	AuthorID          string
	Status            Status
	AssignedReviewers []string
	NeedMoreReviewers bool
	CreatedAt         time.Time
	MergedAt          time.Time
}

func (s Status) String() string {
	return string(s)
}
