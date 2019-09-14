package main

import (
	"net/http"
	"testing"

	"github.com/jcorry/nmap-scan-api/pkg/models/sqlite"
)

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
