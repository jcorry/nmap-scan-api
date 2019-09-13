package models_test

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/jcorry/nmap-scan-api/pkg/models"
)

func Test_ParseXMLData(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  error
	}{
		{
			name:     "Parse included XML file",
			filename: testDataFile("testdata/nmap.results.xml"),
			wantErr:  nil,
		},
		{
			name:     "Parse non XML file",
			filename: testDataFile("testdata/nmap.results.nmap"),
			wantErr:  io.EOF,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Read the file
			b, err := ioutil.ReadFile(fmt.Sprintf("%s", tt.filename))
			if err != nil {
				t.Fatal(err)
			}
			// Hash the file
			hash := sha256.Sum256(b)

			// Parse the bytes
			_, err = models.ParseXMLData(fmt.Sprintf("%x", hash), b)
			// Check for errs
			if err != tt.wantErr {
				t.Log(fmt.Sprintf("Want err %s, Got err %s", tt.wantErr, err))
			}
		})
	}
}

func testDataFile(filename string) string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Error: Couldn't determine working directory: " + err.Error())
	}
	return fmt.Sprintf("%s/%s", wd, filename)
}
