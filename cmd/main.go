package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	th "github.com/LeoUraltsev/PRReviewerService/internal/http/handler/team"
	uh "github.com/LeoUraltsev/PRReviewerService/internal/http/handler/user"
	appmw "github.com/LeoUraltsev/PRReviewerService/internal/http/middleware"
	ts "github.com/LeoUraltsev/PRReviewerService/internal/service/team"
	us "github.com/LeoUraltsev/PRReviewerService/internal/service/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(appmw.ContentTypeApplicationJson)

	teamService := ts.NewService()
	userService := us.NewService()

	teamHandler := th.NewHandler(teamService, teamService)
	userHandler := uh.NewHandler(userService, userService)

	r.Route("/team", func(r chi.Router) {
		r.Post("/add", teamHandler.AddingTeam)
		r.Get("/get", teamHandler.GetTeam)
	})
	r.Route("/users", func(r chi.Router) {
		r.Post("/setIsActive", userHandler.SetIsActive)
		r.Get("/getReview", userHandler.GetReview)
	})
	r.Route("/pullRequest", func(r chi.Router) {
		r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("create pull request"))
		})
		r.Post("/merge", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("merge pull request"))
		})
		r.Post("/reassign", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("reassign pull request"))
		})
	})

	server := http.Server{
		Addr:              "localhost:8080",
		Handler:           r,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       5 * time.Second,
	}
	log.Info("starting server")

	if err := server.ListenAndServe(); err != nil {
		slog.Error("failed to start server", err)
	}

}
