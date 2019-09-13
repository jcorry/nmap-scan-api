package mock

import "github.com/jcorry/nmap-scan-api/pkg/models"

type HostRepo struct{}

func (h *HostRepo) BatchInsert(hosts []*models.Host) (err error) {
	return nil
}

func (h *HostRepo) Insert(host *models.Host) (err error) {
	host.ID = 42
	return nil
}
