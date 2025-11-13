package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

/*
POST /team/add
GET  /team/get
*/

type TeamHandler struct {
	router *chi.Mux
}

func NewTeamHandler(router *chi.Mux) *TeamHandler {

	return &TeamHandler{router: router}
}

func (t *TeamHandler) register() {
	t.router.Post("/team/add", t.getTeam)
	t.router.Get("/team/get", t.getTeam)
}

func (t *TeamHandler) getTeam(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("team get"))
}

func (t *TeamHandler) createTeam(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}
