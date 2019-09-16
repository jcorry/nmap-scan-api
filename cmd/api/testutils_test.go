package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"

	"github.com/jcorry/nmap-scan-api/pkg/models/mock"
)

type testServer struct {
	*httptest.Server
}

func newTestApplication(t *testing.T) *application {
	return &application{
		tmplDir:    os.Getenv("TMPL_DIR"),
		errorLog:   log.New(ioutil.Discard, "", 0),
		infoLog:    log.New(ioutil.Discard, "", 0),
		hostRepo:   &mock.HostRepo{},
		importRepo: &mock.FileImportRepo{},
	}
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)

	return &testServer{ts}
}

func newTestDB(t *testing.T) (*sql.DB, func()) {
	os.Remove("memory:")
	db, err := sql.Open("sqlite3", "file:memory:?cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	// DB migrations
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		t.Fatal(err)
	}

	// Look for migration files relative to this directory
	filePath, err := os.Getwd()
	migrationsPath := fmt.Sprintf("%s/../../sql/migrations", filePath)
	if err != nil {
		t.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://"+migrationsPath),
		"sqlite3", driver)
	if err != nil {
		t.Fatal(err)
	}

	err = m.Up()

	if err != nil {
		t.Fatal(err)
	}

	return db, func() {
		err := m.Down()
		if err != nil {
			t.Fatal(err)
		}
		db.Close()
	}
}

func (ts *testServer) fileRequest(t *testing.T, urlPath string, fileFieldName, filename string) (int, http.Header, []byte) {
	u, err := url.Parse(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	file, err := os.Open(testDataFile(filename))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	reqBody := &bytes.Buffer{}
	writer := multipart.NewWriter(reqBody)
	part, err := writer.CreateFormFile(fileFieldName, filepath.Base(filename))
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatal(err)
	}

	err = writer.Close()
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", u.String(), reqBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rs, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	rsBody, err := ioutil.ReadAll(rs.Body)

	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, rsBody
}

func (ts *testServer) request(t *testing.T, method, urlPath, contentType string, reqBody io.Reader) (int, http.Header, []byte) {
	u, err := url.Parse(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	req := &http.Request{
		Method: "GET",
		URL:    u,
		Body:   ioutil.NopCloser(reqBody),
		Header: map[string][]string{
			"Content-Type": {contentType},
		},
	}

	switch method {
	case http.MethodGet:
		req.Method = http.MethodGet
	case http.MethodPost:
		req.Method = http.MethodPost
	case http.MethodPatch:
		req.Method = http.MethodPatch
	case http.MethodPut:
		req.Method = http.MethodPut
	case http.MethodDelete:
		req.Method = http.MethodDelete
	}

	rs, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	body, err := ioutil.ReadAll(rs.Body)

	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, body
}

func testDataFile(filename string) string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Error: Couldn't determine working directory: " + err.Error())
	}
	return fmt.Sprintf("%s/%s", wd, filename)
}
