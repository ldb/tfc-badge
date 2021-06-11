package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	serverReadTimeoutSeconds  = 60
	serverWriteTimeoutSeconds = 60
)

type AppServer struct {
	Addr        string
	Store       *RunCache
	server      *http.Server
	initialized bool
}

func (a *AppServer) Run() error {
	if !a.initialized {
		a.init()
	}
	return a.server.ListenAndServe()
}

func (a *AppServer) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}

func (a *AppServer) init() {
	a.server = &http.Server{
		Addr:         a.Addr,
		WriteTimeout: serverWriteTimeoutSeconds * time.Second,
		ReadTimeout:  serverReadTimeoutSeconds * time.Second,
		Handler:      a.routes(),
	}

	a.initialized = true
}

func (a *AppServer) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", a.handleIndex())
	mux.HandleFunc("/badge/", a.handleBadge())
	mux.HandleFunc("/run", a.handleRun())

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("request: %s: %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		mux.ServeHTTP(w, r)
	})
}

// Dumb Handler, can be used for readiness and liveliness checks.
func (a *AppServer) handleIndex() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	}
}

func (a *AppServer) handleRun() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		defer request.Body.Close()
		run := new(Run)
		if err := json.NewDecoder(request.Body).Decode(run); err != nil {
			log.Printf("error decoding request: %var", err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		// TODO: Implement HMAC verification.
		a.Store.Set(run.WorkspaceID, run)
		log.Printf("stored new run for workspace %s", run.WorkspaceID)
		writer.WriteHeader(http.StatusOK)
	}
}

func (a *AppServer) handleBadge() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		writer.Header().Set("Content-Type", "image/svg+xml")
		workspaceID := strings.TrimPrefix(request.URL.Path, "/badge/")
		var badge Badge
		run, err := a.Store.Get(workspaceID)
		if err != nil {
			badge = DefaultBadge
			log.Printf("error finding run for workspace %s: %v", workspaceID, err)
		} else {
			badge.FromRun(run)
		}
		if err := badge.Render(writer); err != nil {
			log.Printf("error rendering badge: %v", err)
		}
	}
}
