package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

type AppServer struct {
	Host        string
	Store       *RunCache
	server      *http.Server
	initialized bool
}

func (a *AppServer) Run() {
	if !a.initialized {
		a.init()
	}
	a.server.ListenAndServe()
}

func (a *AppServer) init() {
	a.server = &http.Server{
		Addr:           a.Host,
		ReadTimeout:    1 * time.Second,
		WriteTimeout:   1 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        a.routes(),
	}

	a.initialized = true
}

func (a *AppServer) routes() http.Handler {
	r := http.NewServeMux()

	r.HandleFunc("/", a.handleIndex())
	r.HandleFunc("/badge/", a.handleBadge())
	r.HandleFunc("/run", a.handleRun())
	return r
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
		badge.Render(writer)
	}
}
