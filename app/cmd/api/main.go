package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/uptime-monitor/internal/config"
	"github.com/example/uptime-monitor/internal/health"
	"github.com/example/uptime-monitor/internal/id"
	"github.com/example/uptime-monitor/internal/logging"
	"github.com/example/uptime-monitor/internal/store"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	repo := store.NewDynamo(ctx, cfg.AWSRegion, cfg.SitesTable, cfg.StatusTable)
	logger := logging.New("api")

	mux := http.NewServeMux()
	health.Register(mux)
	mux.HandleFunc("/api/v1/sites", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleListSites(ctx, repo, w)
		case http.MethodPost:
			handleCreateSite(ctx, repo, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/v1/status/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		id := r.URL.Path[len("/api/v1/status/"):]
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		status, err := repo.LatestStatus(ctx, id)
		if err != nil {
			if err == store.ErrNotFound {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			logger.Printf("fetch status: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, status)
	})

	addr := ":" + os.Getenv("API_PORT")
	if addr == ":" {
		addr = ":8080"
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           loggingMiddleware(logger, mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Printf("starting api on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Printf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}

func handleListSites(ctx context.Context, repo store.Repository, w http.ResponseWriter) {
	sites, err := repo.ListSites(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, sites)
}

func handleCreateSite(ctx context.Context, repo store.Repository, w http.ResponseWriter, r *http.Request) {
	type request struct {
		URL string `json:"url"`
	}
	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.URL == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	site := store.Site{
		ID:        id.New(),
		URL:       req.URL,
		CreatedAt: time.Now().UTC(),
	}
	if err := repo.PutSite(ctx, site); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, site)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func loggingMiddleware(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
