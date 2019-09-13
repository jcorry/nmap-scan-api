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

	"github.com/jcorry/nmap-scan-api/pkg/models/mock"
)

type testServer struct {
	*httptest.Server
}

func newTestApplication(t *testing.T) *application {
	return &application{
		errorLog:   log.New(ioutil.Discard, "", 0),
		infoLog:    log.New(ioutil.Discard, "", 0),
		hostRepo:   &mock.HostRepo{},
		importRepo: &mock.FileImportRepo{},
	}
}

func newTestDB(t *testing.T) *sql.DB {
	sqlite_test.newTestDB(t)
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)

	return &testServer{ts}
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

func (ts *testServer) request(t *testing.T, method string, urlPath string, reqBody io.Reader) (int, http.Header, []byte) {
	u, err := url.Parse(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	req := &http.Request{
		Method: "PATCH",
		URL:    u,
		Body:   ioutil.NopCloser(reqBody),
		Header: map[string][]string{
			"Content-Type": {"application/json"},
		},
	}

	switch method {
	case "get":
		req.Method = "GET"
	case "post":
		req.Method = "POST"
	case "patch":
		req.Method = "PATCH"
	case "put":
		req.Method = "PUT"
	case "delete":
		req.Method = "DELETE"
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
