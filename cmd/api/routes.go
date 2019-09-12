package main

import (
	"net/http"

	"github.com/bmizerany/pat"
)

func (a *application) routes() http.Handler {
	mux := pat.New()
	// API routes
	mux.Post("/api/v1/nmap/upload", http.HandlerFunc(a.upload))

	return a.logRequest(mux)
}
