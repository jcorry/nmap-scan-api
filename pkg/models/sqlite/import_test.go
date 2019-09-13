package sqlite_test

import (
	"strings"
	"testing"

	"github.com/jcorry/nmap-scan-api/pkg/models"

	"github.com/jcorry/nmap-scan-api/pkg/models/sqlite"
)

func Test_FileImportInsert(t *testing.T) {

	db, teardown := newTestDB(t)
	defer teardown()

	tests := []struct {
		name            string
		wantErrContains string
	}{
		{
			name:            "Successful Insert",
			wantErrContains: "",
		},
		{
			name:            "Unique constraint",
			wantErrContains: "UNIQUE",
		},
	}

	fileImportToSave := &models.FileImport{
		FileID: "xxxxx",
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fi := sqlite.FileImportRepo{DB: db}

			err := fi.Insert(fileImportToSave)

			if err != nil || tt.wantErrContains != "" {
				if !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Fatalf("Want err `%s` to contain `%s`", err, tt.wantErrContains)
				}
			}
		})
	}
}
