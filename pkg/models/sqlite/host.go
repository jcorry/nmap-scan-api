package sqlite

import (
	"database/sql"
	"fmt"
	"sort"
	"time"

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
			_, err = insertAddressStmt.Exec(h.ID, a.Addr.String, a.AddrType.String)
			if err != nil {
				return
			}
			a.HostID = models.ToNullInt64(h.ID)
		}

		// Insert hostnames
		for _, hn := range h.Hostnames {
			_, err = insertHostnameStmt.Exec(h.ID, hn.Name.String, hn.Type.String)
			if err != nil {
				return
			}
			hn.HostID = models.ToNullInt64(h.ID)
		}

		// Insert ports
		for _, p := range h.Ports {
			_, err = insertPortStmt.Exec(h.ID, p.Protocol.String, p.PortID.Int64, p.Owner.String, p.Service.String)
			if err != nil {
				return
			}
			p.HostID = models.ToNullInt64(h.ID)
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
		_, err = stmt.Exec(host.ID, a.Addr.String, a.AddrType.String)
		if err != nil {
			return
		}
		a.HostID = models.ToNullInt64(host.ID)
	}

	// Insert hostnames
	for _, hn := range host.Hostnames {
		stmt, err = h.DB.Prepare(`INSERT INTO hostnames (host_id, name, type) VALUES (?, ?, ?)`)
		if err != nil {
			return
		}
		_, err = stmt.Exec(host.ID, hn.Name.String, hn.Type.String)
		if err != nil {
			return
		}
		hn.HostID = models.ToNullInt64(host.ID)
	}

	// Insert ports
	for _, p := range host.Ports {
		stmt, err = h.DB.Prepare(`INSERT INTO ports (host_id, protocol, port_id, owner, service) VALUES (?, ?, ?, ?, ?)`)
		if err != nil {
			return
		}
		_, err = stmt.Exec(host.ID, p.Protocol.String, p.PortID.Int64, p.Owner.String, p.Service.String)
		if err != nil {
			return
		}
		p.HostID = models.ToNullInt64(host.ID)
	}

	return
}

// List retrieves a list of Hosts with each of its subtypes
func (h *HostRepo) List(start, length int) (meta *models.Meta, hosts []*models.Host, err error) {
	var offset string
	if start > 0 {
		offset = fmt.Sprintf(` OFFSET %d`, start)
	}
	var limit string
	if length > 0 {
		limit = fmt.Sprintf(` LIMIT %d`, length)
	}
	stmt, err := h.DB.Prepare(fmt.Sprintf(`SELECT h.id, h.file_id, h.starttime, h.endtime, h.comment, 
													     a.host_id, a.addr, a.addrtype, 
														 p.host_id, p.protocol, p.port_id, p.owner, p.service, 
														 hn.host_id, hn.name, hn.type
											        FROM (SELECT * FROM hosts%s%s) h
											   LEFT JOIN addresses a ON a.host_id = h.id
											   LEFT JOIN ports p ON a.host_id = p.host_id 
											   LEFT JOIN hostnames hn on a.host_id = hn.host_id 
												ORDER BY h.id DESC`, limit, offset))
	if err != nil {
		return nil, nil, err
	}
	rows, err := stmt.Query()

	var keys []int
	hostMap := make(map[int]*models.Host)
	addrMap := make(map[int][]*models.Address)
	portMap := make(map[int][]*models.Port)
	hostnameMap := make(map[int][]*models.Hostname)

	for rows.Next() {

		h := &models.Host{
			ID:        0,
			FileID:    "",
			StartTime: time.Time{},
			EndTime:   time.Time{},
			Comment:   "",
			Status:    "",
			Hostnames: nil,
			Addresses: nil,
			Ports:     nil,
		}

		a := &models.Address{}

		p := &models.Port{}

		hn := &models.Hostname{}

		err = rows.Scan(&h.ID, &h.FileID, &h.StartTime, &h.EndTime, &h.Comment, &a.HostID, &a.Addr, &a.AddrType, &p.HostID, &p.Protocol, &p.PortID, &p.Owner, &p.Service, &hn.HostID, &hn.Name, &hn.Type)
		if err != nil {
			return
		}
		if a.HostID == models.ToNullInt64(h.ID) {
			addrMap[h.ID] = append(addrMap[h.ID], a)
		}

		if p.HostID == models.ToNullInt64(h.ID) {
			portMap[h.ID] = append(portMap[h.ID], p)
		}

		if hn.HostID == models.ToNullInt64(h.ID) {
			hostnameMap[h.ID] = append(hostnameMap[h.ID], hn)
		}

		hostMap[h.ID] = h
	}

	for _, host := range hostMap {
		host.Addresses = addrMap[host.ID]
		host.Ports = portMap[host.ID]
		host.Hostnames = hostnameMap[host.ID]

		keys = append(keys, host.ID)
	}

	// Sort the map by ID
	sort.Ints(keys)
	for _, k := range keys {
		hosts = append(hosts, hostMap[k])
	}

	total, err := h.Count()
	if err != nil {
		return
	}

	meta = &models.Meta{
		Start:  start,
		Length: length,
		Total:  total,
	}

	return
}

// Count returns a count of the hosts in the DB with addresses matching `addr`
func (h *HostRepo) Count() (count int, err error) {
	stmt, err := h.DB.Prepare(`SELECT COUNT(*) FROM hosts`)
	if err != nil {
		return
	}
	err = stmt.QueryRow().Scan(&count)
	if err != nil {
		return
	}
	return count, nil
}
