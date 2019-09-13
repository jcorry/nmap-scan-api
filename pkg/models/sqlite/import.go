package sqlite

import (
	"database/sql"
	"time"

	"github.com/jcorry/nmap-scan-api/pkg/models"
)

type FileImportRepo struct {
	DB *sql.DB
}

func (i *FileImportRepo) Insert(fileImport *models.FileImport) (err error) {
	var stmt *sql.Stmt
	stmt, err = i.DB.Prepare(`INSERT INTO imports (file_id, created) VALUES (?, ?)`)
	if err != nil {
		return
	}

	var res sql.Result
	created := time.Now()

	res, err = stmt.Exec(fileImport.FileID, time.Now())
	if err != nil {
		return
	}

	importID, err := res.LastInsertId()
	if err != nil {
		return
	}

	fileImport.ID = int(importID)
	fileImport.Created = created

	return
}
