package team

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/LeoUraltsev/PRReviewerService/internal/domain"
	"github.com/go-chi/render"
)

/*
POST /team/add
GET  /team/get
*/

var (
	notFound      = "NOT_FOUND"
	teamExists    = "TEAM_EXISTS"
	incorrectData = "INCORRECT_DATA"
)

type Saver interface {
	Save(ctx context.Context, team *domain.Team) error
}

type teamMember struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}
type team struct {
	TeamName string       `json:"team_name"`
	Members  []teamMember `json:"members"`
}

type responseAddTeam struct {
	Team  team        `json:"team,omitempty"`
	Error responseErr `json:"error,omitempty"`
}

type responseErr struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Handler struct {
	saver Saver
}

func NewHandler(saver Saver) *Handler {
	return &Handler{
		saver: saver,
	}
}

func (h *Handler) AddingTeam(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var t team
	err := render.DecodeJSON(r.Body, &t)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		render.JSON(w, r, responseAddTeam{
			Error: responseErr{
				Code:    incorrectData,
				Message: "failed getting body",
			},
		})

		return
	}
	defer func() {
		_ = r.Body.Close()
	}()

	teamDomain := toDomain(t)
	err = h.saver.Save(context.Background(), teamDomain)
	if err != nil {
		if errors.Is(err, domain.ErrTeamExists) {
			w.WriteHeader(http.StatusBadRequest)

			render.JSON(w, r, responseAddTeam{
				Error: responseErr{
					Code:    teamExists,
					Message: fmt.Sprintf("%s already exists", t.TeamName),
				},
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func (h *Handler) GetTeam(w http.ResponseWriter, r *http.Request) {

}

// todo: добавить обработку ошибок перед сдачей: deadline - 16.11.2025
func toDomain(team team) *domain.Team {
	m := make([]domain.User, len(team.Members))

	for i, member := range team.Members {
		m[i] = domain.User{
			UserID:   member.UserId,
			Username: member.Username,
			TeamName: team.TeamName,
			IsActive: member.IsActive,
		}
	}

	return &domain.Team{
		TeamName: team.TeamName,
		Members:  m,
	}
}
