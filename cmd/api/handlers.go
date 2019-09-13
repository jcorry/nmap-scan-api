package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/jcorry/nmap-scan-api/pkg/models"
)

func (a *application) upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		a.serverError(w, err)
		return
	}

	file, _, err := r.FormFile("file")
	defer file.Close()

	if err != nil {
		a.serverError(w, err)
		return
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		a.serverError(w, err)
		return
	}

	// Get file hash
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		a.serverError(w, err)
		return
	}

	data, err := models.ParseXMLData(fmt.Sprintf("%x", h.Sum(nil)), b)
	if err != nil {
		a.serverError(w, err)
		return
	}

	// w.WriteHeader(http.StatusCreated)
	a.jsonResponse(w, data)
}
