package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/LeoUraltsev/PRReviewerService/internal/config"
	"github.com/LeoUraltsev/PRReviewerService/internal/http/handler/pull_request"
	th "github.com/LeoUraltsev/PRReviewerService/internal/http/handler/team"
	uh "github.com/LeoUraltsev/PRReviewerService/internal/http/handler/user"
	appmw "github.com/LeoUraltsev/PRReviewerService/internal/http/middleware"
	pr "github.com/LeoUraltsev/PRReviewerService/internal/service/pull_request"
	ts "github.com/LeoUraltsev/PRReviewerService/internal/service/team"
	us "github.com/LeoUraltsev/PRReviewerService/internal/service/user"
	"github.com/LeoUraltsev/PRReviewerService/internal/storage/pg"
	pullStorage "github.com/LeoUraltsev/PRReviewerService/internal/storage/pg/pull_request"
	teamStorage "github.com/LeoUraltsev/PRReviewerService/internal/storage/pg/team"
	userStorage "github.com/LeoUraltsev/PRReviewerService/internal/storage/pg/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.LogLever,
	}))

	ctx := context.Background()

	s, err := pg.NewStorage(ctx, log, cfg)
	defer s.Close(ctx)
	if err != nil {
		fmt.Printf("Stopping application: %v\n", err)
		os.Exit(1)
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(appmw.ContentTypeApplicationJson)

	prStorage := pullStorage.NewStorage(log, s.Pool)
	tStorage := teamStorage.NewStorage(log, s.Pool)
	uStorage := userStorage.NewStorage(log, s.Pool)

	teamService := ts.NewService(uStorage, tStorage)
	userService := us.NewService(prStorage, uStorage)
	prService := pr.NewService(prStorage, uStorage)

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
		Addr:              fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Handler:           r,
		ReadTimeout:       cfg.ReadTimeout * time.Second,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout * time.Second,
		WriteTimeout:      cfg.WriteTimeout * time.Second,
		IdleTimeout:       cfg.IdleTimeout * time.Second,
	}
	log.Info("starting server")

	if err = server.ListenAndServe(); err != nil {
		slog.Error("failed to start server", slog.String("error", err.Error()))
	}

}
