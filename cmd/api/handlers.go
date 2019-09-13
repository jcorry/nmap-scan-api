package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/jcorry/nmap-scan-api/pkg/models"
)

func (a *application) upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		a.clientError(w, err, 400)
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
	fileID := fmt.Sprintf("%x", h.Sum(nil))

	hosts, err := models.ParseXMLData(fileID, b)
	if err != nil {
		a.serverError(w, err)
		return
	}
	// Log the import
	fileImport := &models.FileImport{
		FileID: fileID,
	}
	err = a.importRepo.Insert(fileImport)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			a.clientError(w, errors.New("That file has already been imported"), 400)
		} else {
			a.serverError(w, err)
		}
		return
	}

	err = a.hostRepo.BatchInsert(hosts)
	if err != nil {
		a.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte{})
	return
}
