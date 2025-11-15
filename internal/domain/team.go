package domain

import "errors"

var (
	ErrTeamNotFound = errors.New("team not found")
	ErrTeamExists   = errors.New("team already exists")
)

type Team struct {
	TeamName string
	Members  []*User
}
