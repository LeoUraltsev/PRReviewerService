package domain

import (
	"errors"
	"time"
)

var (
	ErrPRNotFound       = errors.New("pr not found")
	ErrPRAlreadyExists  = errors.New("pr already exists")
	ErrReassignPRMerged = errors.New("cannot reassign on merged PR")
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

func (p *PullRequest) ChangeStatus(status Status) {
	p.Status = status
}
