package main

import (
	"bytes"
	"net/http"
	"strings"
	"testing"

	"github.com/jcorry/nmap-scan-api/pkg/models/sqlite"
)

func Test_List(t *testing.T) {
	app := newTestApplication(t)
	testDB, cleanup := newTestDB(t)
	defer cleanup()
	app.hostRepo = &sqlite.HostRepo{DB: testDB}
	app.importRepo = &sqlite.FileImportRepo{DB: testDB}

	// Insert data
	dataFilePath := "../../pkg/models/testdata/nmap.results.xml"
	urlPath := "/api/v1/nmap/upload"

	ts := newTestServer(t, app.routes())
	defer ts.Close()
	// Easiest way to populate the DB is to just upload the file ¯\_(ツ)_/¯
	code, _, _ := ts.fileRequest(t, urlPath, "file", dataFilePath)
	if code != http.StatusCreated {
		t.Fatalf("Failed to upload XML file to test server")
	}

	tests := []struct {
		name        string
		url         string
		contentType string
		start       int
		length      int
		wantCode    int
		wantBody    string
	}{
		{
			"Success",
			"/nmap/list",
			"text/html",
			0,
			0,
			200,
			"Host",
		},
		{
			"Success JSON",
			"/api/v1/nmap",
			"application/json",
			0,
			0,
			200,
			"items",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.request(t, http.MethodGet, tt.url, tt.contentType, bytes.NewReader(nil))
			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !strings.Contains(string(body), tt.wantBody) {
				t.Errorf("Wanted body to contain the string %s, not found", tt.wantBody)
			}
		})
	}
}

func Test_Upload(t *testing.T) {
	app := newTestApplication(t)
	testDB, cleanup := newTestDB(t)
	defer cleanup()
	app.hostRepo = &sqlite.HostRepo{DB: testDB}
	app.importRepo = &sqlite.FileImportRepo{DB: testDB}

	dataFilePath := "../../pkg/models/testdata/nmap.results.xml"

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		wantCode int
	}{
		{
			name:     "200 Passing",
			wantCode: http.StatusCreated,
		},
		{
			name:     "400 Unique constraint fails on 2nd run",
			wantCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlPath := "/api/v1/nmap/upload"

			code, _, body := ts.fileRequest(t, urlPath, "file", dataFilePath)
			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
				t.Errorf("%s", body)
			}
		})
	}
}
