package domain

import "errors"

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrIncorrectAdminToken = errors.New("incorrect admin token")
)

type User struct {
	UserID   string
	Username string
	TeamName string
	IsActive bool
}

func (u *User) ChangeActive(active bool) {
	u.IsActive = active
}
