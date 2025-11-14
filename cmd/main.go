package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	appmw "github.com/LeoUraltsev/PRReviewerService/internal/http/middleware"
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

	r.Route("/team", func(r chi.Router) {
		r.Post("/add", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("add team"))
		})
		r.Get("/get", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("get team"))
		})
	})
	r.Route("/users", func(r chi.Router) {
		r.Post("/setIsActive", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("set isActive"))
		})
		r.Post("/getReview", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("get review"))
		})
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
