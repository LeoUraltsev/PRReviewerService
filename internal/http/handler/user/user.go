package user

import (
	"context"
	"errors"
	"net/http"

	"github.com/LeoUraltsev/PRReviewerService/internal/domain"
	e "github.com/LeoUraltsev/PRReviewerService/internal/http/handler/helper/err"
	"github.com/go-chi/render"
)

type Updater interface {
	UpdateIsActive(ctx context.Context, userId string, isActive bool) (*domain.User, error)
}

type Getter interface {
	GetUserPullRequest(ctx context.Context, userId string) ([]*domain.PullRequest, error)
}

type isActiveRequest struct {
	UserId   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type isActiveResponse struct {
	user `json:"user"`
}

type getReviewResponse struct {
	UserId       string              `json:"user_id"`
	PullRequests []smallPullRequests `json:"pull_requests"`
}

type smallPullRequests struct {
	PullRequestId   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorId        string `json:"author_id"`
	Status          string `json:"status"`
}

type user struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type Handler struct {
	updater Updater
	getter  Getter
}

func NewHandler(updater Updater, getter Getter) *Handler {
	return &Handler{
		updater: updater,
		getter:  getter,
	}
}

func (h *Handler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req isActiveRequest
	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, e.IncorrectDataError())
		return
	}

	u, err := h.updater.UpdateIsActive(r.Context(), req.UserId, req.IsActive)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) || errors.Is(err, domain.ErrIncorrectAdminToken) {
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, e.NotFoundError())
		}

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, e.InternalServerError())
		return
	}

	resp := isActiveResponse{
		user: userDomainTo(u),
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, resp)
}

func (h *Handler) GetReview(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, e.IncorrectDataError())
		return
	}

	prDomain, err := h.getter.GetUserPullRequest(r.Context(), userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, e.NotFoundError())
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, e.InternalServerError())
		return
	}

	pr := make([]smallPullRequests, len(prDomain))

	for i, v := range prDomain {
		pr[i] = prDomainToSmallPullRequests(v)
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, getReviewResponse{
		UserId:       userID,
		PullRequests: pr,
	})

}

func userDomainTo(u *domain.User) user {
	return user{
		UserId:   u.UserID,
		Username: u.Username,
		TeamName: u.TeamName,
		IsActive: u.IsActive,
	}
}

func prDomainToSmallPullRequests(pr *domain.PullRequest) smallPullRequests {
	return smallPullRequests{
		PullRequestId:   pr.ID,
		PullRequestName: pr.Name,
		AuthorId:        pr.AuthorID,
		Status:          pr.Status.String(),
	}
}
