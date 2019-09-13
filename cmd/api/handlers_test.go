package main

import (
	"net/http"
	"testing"
)

func Test_Upload(t *testing.T) {
	app := newTestApplication(t)
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlPath := "/api/v1/nmap/upload"

			code, _, body := ts.fileRequest(t, urlPath, "file", "../../pkg/models/testdata/nmap.results.xml")
			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
				t.Errorf("%s", body)
			}
		})
	}
}
