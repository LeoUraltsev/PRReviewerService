package user

import (
	"context"
	"net/http"

	"github.com/LeoUraltsev/PRReviewerService/internal/domain"
	"github.com/LeoUraltsev/PRReviewerService/internal/http/handler/helper"
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

func NewHandler(updater Updater) *Handler {
	return &Handler{
		updater: updater,
	}
}

// todo: обработка всех ошибок
func (h *Handler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req isActiveRequest
	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, helper.NewErrorResponse("INCORRECT_DATA", "failed getting body"))
	}

	u, err := h.updater.UpdateIsActive(r.Context(), req.UserId, req.IsActive)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, helper.NewErrorResponse("INTERNAL_ERROR", "internal server error"))
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
		render.JSON(w, r, helper.NewErrorResponse("NO_USER_ID", "no user id"))
	}

	prDomain, err := h.getter.GetUserPullRequest(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, helper.NewErrorResponse("INTERNAL_ERROR", "internal server error"))
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
