package main

import (
	"net/http"
)

func (a *application) upload(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusCreated)
}
