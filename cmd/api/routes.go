package main

import (
	"net/http"

	"github.com/bmizerany/pat"
)

func (a *application) routes() http.Handler {
	mux := pat.New()
	// API routes
	mux.Post("/api/v1/nmap/upload", http.HandlerFunc(a.upload))

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	mux.Get("/nmap/list", http.HandlerFunc(a.list))

	return a.logRequest(mux)
}
