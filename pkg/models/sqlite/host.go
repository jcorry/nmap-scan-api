package sqlite

import (
	"database/sql"

	"github.com/jcorry/nmap-scan-api/pkg/models"
)

// HostRepo is a repository struct to receive functions interacting with Host data storage/retrieval
type HostRepo struct {
	DB *sql.DB
}

// BatchInsert inserts all of a collection of Hosts
func (h *HostRepo) BatchInsert(hosts []*models.Host) (err error) {
	insertHostStmt, err := h.DB.Prepare(`INSERT INTO hosts (file_id, starttime, endtime, comment) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return
	}
	insertAddressStmt, err := h.DB.Prepare(`INSERT INTO addresses (host_id, addr, addrtype) VALUES (?, ?, ?)`)
	if err != nil {
		return
	}
	insertHostnameStmt, err := h.DB.Prepare(`INSERT INTO hostnames (host_id, name, type) VALUES (?, ?, ?)`)
	if err != nil {
		return
	}
	insertPortStmt, err := h.DB.Prepare(`INSERT INTO ports (host_id, protocol, port_id, owner, service) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return
	}

	tx, err := h.DB.Begin()
	defer tx.Rollback()

	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	for _, h := range hosts {
		// Insert the host
		var res sql.Result
		res, err = insertHostStmt.Exec(h.FileID, h.StartTime, h.EndTime, h.Comment)
		if err != nil {
			return
		}
		var hostID int64
		hostID, err = res.LastInsertId()
		if err != nil {
			return
		}

		h.ID = int(hostID)

		// Insert Addresses
		for _, a := range h.Addresses {
			_, err = insertAddressStmt.Exec(h.ID, a.Addr, a.AddrType)
			if err != nil {
				return
			}
			a.HostID = h.ID
		}

		// Insert hostnames
		for _, hn := range h.Hostnames {
			_, err = insertHostnameStmt.Exec(h.ID, hn.Name, hn.Type)
			if err != nil {
				return
			}
			hn.HostID = h.ID
		}

		// Insert ports
		for _, p := range h.Ports {
			_, err = insertPortStmt.Exec(h.ID, p.Protocol, p.PortID, p.Owner, p.Service)
			if err != nil {
				return
			}
			p.HostID = h.ID
		}
	}
	return
}

// Insert adds a Host with each of its sub types to the DB
func (h *HostRepo) Insert(host *models.Host) (err error) {
	// Save the host first, get the ID
	tx, err := h.DB.Begin()
	defer tx.Rollback()

	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Insert host and get it's ID
	stmt, err := h.DB.Prepare(`INSERT INTO hosts (file_id, starttime, endtime, comment) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return
	}

	res, err := stmt.Exec(host.FileID, host.StartTime, host.EndTime, host.Comment)
	if err != nil {
		return
	}

	hostID, err := res.LastInsertId()
	if err != nil {
		return
	}

	host.ID = int(hostID)
	// Insert addresses
	for _, a := range host.Addresses {
		stmt, err = h.DB.Prepare(`INSERT INTO addresses (host_id, addr, addrtype) VALUES (?, ?, ?)`)
		if err != nil {
			return
		}
		_, err = stmt.Exec(host.ID, a.Addr, a.AddrType)
		if err != nil {
			return
		}
		a.HostID = host.ID
	}

	// Insert hostnames
	for _, hn := range host.Hostnames {
		stmt, err = h.DB.Prepare(`INSERT INTO hostnames (host_id, name, type) VALUES (?, ?, ?)`)
		if err != nil {
			return
		}
		_, err = stmt.Exec(host.ID, hn.Name, hn.Type)
		if err != nil {
			return
		}
		hn.HostID = host.ID
	}

	// Insert ports
	for _, p := range host.Ports {
		stmt, err = h.DB.Prepare(`INSERT INTO ports (host_id, protocol, port_id, owner, service) VALUES (?, ?, ?, ?, ?)`)
		if err != nil {
			return
		}
		_, err = stmt.Exec(host.ID, p.Protocol, p.PortID, p.Owner, p.Service)
		if err != nil {
			return
		}
		p.HostID = host.ID
	}

	return
}
