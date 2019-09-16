package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/jcorry/nmap-scan-api/pkg/models"
)

// Handler for API file upload function
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

	// ReadAll is a red flag for resource use, could this be a io.Reader instead? Ultimately, I have to parse
	// the bytes in the file and can't do that unless I get the bytes. A better architecture would be to separate this.
	// Remove the file parsing to a separate service that can work async. Our API would simply accept file uploads, put
	// the files in an S3 bucket and background workers could handle parsing/ingesting them.
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

// Handler for HTML list view
func (a *application) list(w http.ResponseWriter, r *http.Request) {
	// Get addr from query string
	q := r.URL.Query()
	// Get the data
	var count = 0
	count, _ = a.hostRepo.Count()

	// Get the files from wd so tests can find them too
	templates := []string{
		fmt.Sprintf("%s/%s", a.tmplDir, "layout.tmpl"),
		fmt.Sprintf("%s/%s", a.tmplDir, "list.tmpl"),
	}

	type Data struct {
		Title string
		Meta  *models.Meta
		Hosts []*models.Host
	}

	start, err := strconv.Atoi(q.Get("start"))
	length, err := strconv.Atoi(q.Get("length"))
	if err != nil {
		// If that didn't work just default them
		start = 0
		length = 20
	}
	meta, hosts, err := a.hostRepo.List(start, length)

	if err != nil {
		a.serverError(w, err)
		return
	}

	meta.Total = count

	data := Data{
		Title: "nmap list",
		Meta:  meta,
		Hosts: hosts,
	}

	// JSON return allows the same handler to be used for API endpoint
	if r.Header.Get("Content-Type") == "application/json" {
		type JSONData struct {
			Meta  *models.Meta   `json:"meta"`
			Hosts []*models.Host `json:"items"`
		}

		jsonData := JSONData{
			Meta:  meta,
			Hosts: hosts,
		}
		a.jsonResponse(w, jsonData)
		return
	}

	// If we didn't return JSON, we'll return HTML
	t, err := template.ParseFiles(templates...)
	if err != nil {
		fmt.Println(err)
		a.errorLog.Output(2, err.Error())
	}

	if err := t.Execute(w, data); err != nil {
		a.clientError(w, err, http.StatusInternalServerError)
		return
	}
}
