package mock

import (
	"time"

	"github.com/jcorry/nmap-scan-api/pkg/models"
)

type FileImportRepo struct{}

func (f *FileImportRepo) Insert(fileImport *models.FileImport) (err error) {
	fileImport.Created = time.Now()
	fileImport.ID = 42
	return nil
}
