package team

import (
	"context"
	"errors"
	"net/http"

	"github.com/LeoUraltsev/PRReviewerService/internal/domain"
	e "github.com/LeoUraltsev/PRReviewerService/internal/http/handler/helper/err"
	"github.com/go-chi/render"
)

type Saver interface {
	Save(ctx context.Context, team *domain.Team) error
}

type Getter interface {
	Get(ctx context.Context, teamName string) (*domain.Team, error)
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
	Team team `json:"team"`
}

type Handler struct {
	saver  Saver
	getter Getter
}

func NewHandler(saver Saver, getter Getter) *Handler {
	return &Handler{
		saver:  saver,
		getter: getter,
	}
}

func (h *Handler) AddingTeam(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var t team

	err := render.DecodeJSON(r.Body, &t)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, e.IncorrectDataError())
		return
	}

	teamDomain := toDomain(t)
	err = h.saver.Save(r.Context(), teamDomain)
	if err != nil {
		if errors.Is(err, domain.ErrTeamExists) {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, e.TeamExistsError())
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, e.InternalServerError())
		return
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, responseAddTeam{
		Team: t,
	})
}

func (h *Handler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		w.WriteHeader(http.StatusNotFound)
		render.JSON(w, r, e.NotFoundError())
		return
	}

	teamDomain, err := h.getter.Get(r.Context(), teamName)
	if err != nil {
		if errors.Is(err, domain.ErrTeamNotFound) {
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, e.NotFoundError())
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, e.InternalServerError())
		return
	}

	t := domainTo(teamDomain)

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, responseAddTeam{
		Team: t,
	})
}

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

func domainTo(t *domain.Team) team {
	if t == nil {
		return team{}
	}
	members := make([]teamMember, len(t.Members))
	for i, member := range t.Members {
		members[i] = teamMember{
			UserId:   member.UserID,
			Username: member.Username,
			IsActive: member.IsActive,
		}
	}

	return team{
		TeamName: t.TeamName,
		Members:  members,
	}

}
