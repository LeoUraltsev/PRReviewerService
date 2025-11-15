package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/LeoUraltsev/PRReviewerService/internal/http/handler/pull_request"
	th "github.com/LeoUraltsev/PRReviewerService/internal/http/handler/team"
	uh "github.com/LeoUraltsev/PRReviewerService/internal/http/handler/user"
	appmw "github.com/LeoUraltsev/PRReviewerService/internal/http/middleware"
	pr "github.com/LeoUraltsev/PRReviewerService/internal/service/pull_request"
	ts "github.com/LeoUraltsev/PRReviewerService/internal/service/team"
	us "github.com/LeoUraltsev/PRReviewerService/internal/service/user"
	"github.com/LeoUraltsev/PRReviewerService/internal/storage/pg"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	ctx := context.Background()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(appmw.ContentTypeApplicationJson)

	s, err := pg.NewStorage(ctx, log)
	if err != nil {
		fmt.Printf("Stopping application: %v\n", err)
		os.Exit(1)
	}
	defer s.Close(ctx)

	teamService := ts.NewService()
	userService := us.NewService()
	prService := pr.NewService()

	teamHandler := th.NewHandler(teamService, teamService)
	userHandler := uh.NewHandler(userService, userService)
	prHandler := pull_request.NewHandler(prService, prService)

	r.Route("/team", func(r chi.Router) {
		r.Post("/add", teamHandler.AddingTeam)
		r.Get("/get", teamHandler.GetTeam)
	})
	r.Route("/users", func(r chi.Router) {
		r.Post("/setIsActive", userHandler.SetIsActive)
		r.Get("/getReview", userHandler.GetReview)
	})
	r.Route("/pullRequest", func(r chi.Router) {
		r.Post("/create", prHandler.CreatePullRequest)
		r.Post("/merge", prHandler.MergePullRequest)
		r.Post("/reassign", prHandler.ReassignPullRequest)
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
