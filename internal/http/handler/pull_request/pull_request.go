package pull_request

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/LeoUraltsev/PRReviewerService/internal/domain"
	e "github.com/LeoUraltsev/PRReviewerService/internal/http/handler/helper/err"
	"github.com/go-chi/render"
)

type Saver interface {
	SavePullRequest(ctx context.Context, prID string, prName string, authorID string) (*domain.PullRequest, error)
}

type Updater interface {
	MergePullRequest(ctx context.Context, prID string) (*domain.PullRequest, error)
	ReassignReviewerPullRequest(ctx context.Context, prID string, reviewerID string) (*domain.PullRequest, error)
}

type createPRRequest struct {
	PullRequestId   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorId        string `json:"author_id"`
}

type createPRResponse struct {
	PullRequest pullRequest `json:"pr"`
}

type pullRequest struct {
	PullRequestId     string   `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	AuthorId          string   `json:"author_id"`
	Status            string   `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
}

type mergedPRRequest struct {
	PullRequestId string `json:"pull_request_id"`
}

type mergedPRResponse struct {
	PullRequest pullRequest `json:"pr"`
	MergedAt    time.Time   `json:"mergedAt"`
}

type reassignPRRequest struct {
	PullRequestId string `json:"pull_request_id"`
	OldReviewerId string `json:"old_reviewer_id"`
}

type reassignPRResponse struct {
	PullRequest pullRequest `json:"pr"`
	ReplacedBy  string      `json:"replaced_by"`
}

type Handler struct {
	saver   Saver
	updater Updater
}

func New(saver Saver, updater Updater) *Handler {
	return &Handler{
		saver:   saver,
		updater: updater,
	}
}

func (h *Handler) CreatePullRequest(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req createPRRequest
	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, e.IncorrectDataError())
		return
	}

	prDomain, err := h.saver.SavePullRequest(r.Context(), req.PullRequestId, req.PullRequestName, req.AuthorId)
	if err != nil {
		if errors.Is(err, domain.ErrPRNotFound) {
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, e.NotFoundError())
			return
		}

		if errors.Is(err, domain.ErrPRAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			render.JSON(w, r, e.PRExistsError())
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, e.InternalServerError())
		return
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, createPRResponse{
		PullRequest: domainToPullRequest(prDomain),
	})

}

func (h *Handler) MergePullRequest(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req mergedPRRequest
	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, e.IncorrectDataError())
		return
	}

	prDomain, err := h.updater.MergePullRequest(r.Context(), req.PullRequestId)
	if err != nil {
		if errors.Is(err, domain.ErrPRNotFound) {
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, e.NotFoundError())
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	mergedAt := prDomain.MergedAt
	pr := domainToPullRequest(prDomain)
	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, mergedPRResponse{
		PullRequest: pr,
		MergedAt:    mergedAt,
	})
}

func (h *Handler) ReassignPullRequest(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req reassignPRRequest
	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, e.IncorrectDataError())
		return
	}

	prDomain, err := h.updater.ReassignReviewerPullRequest(r.Context(), req.PullRequestId, req.OldReviewerId)
	if err != nil {
		if errors.Is(err, domain.ErrPRNotFound) {
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, e.NotFoundError())
			return
		}

		if errors.Is(err, domain.ErrReassignPRMerged) {
			w.WriteHeader(http.StatusConflict)
			render.JSON(w, r, e.PRMergedError())
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, e.InternalServerError())
		return
	}

	pr := domainToPullRequest(prDomain)
	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, reassignPRResponse{
		PullRequest: pr,
		ReplacedBy:  req.OldReviewerId,
	})
}

func domainToPullRequest(pr *domain.PullRequest) pullRequest {
	return pullRequest{
		PullRequestId:     pr.ID,
		PullRequestName:   pr.Name,
		AuthorId:          pr.AuthorID,
		Status:            pr.Status.String(),
		AssignedReviewers: pr.AssignedReviewers,
	}
}
