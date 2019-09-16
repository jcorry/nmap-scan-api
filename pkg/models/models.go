package models

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"time"

	nmap "github.com/tomsteele/go-nmap"
)

// Meta is metadata about a response
type Meta struct {
	Start  int `json:"start"`
	Length int `json:"length"`
	Total  int `json:"total"`
}

// FileImport is the importation of a single nmap file
type FileImport struct {
	ID      int       `json:"id" db:"id"`
	FileID  string    `json:"fileid" db:"file_id"`
	Created time.Time `json:"created" db:"created"`
}

// Host is an internal type that contains a subset of the data in a go-nmap.Host
type Host struct {
	ID        int         `json:"id" db:"id"`
	FileID    string      `json:"fileid" db:"file_id"`
	StartTime time.Time   `json:"starttime" db:"starttime"`
	EndTime   time.Time   `json:"endtime" db:"endtime"`
	Comment   string      `json:"comment" db:"comment"`
	Status    string      `json:"status" db:"status"`
	Hostnames []*Hostname `json:"hostnames"`
	Addresses []*Address  `json:"addresses"`
	Ports     []*Port     `json:"ports"`
}

// Address is an internal type that contains a subset of the data in a go-nmap.Address
type Address struct {
	HostID   NullInt64  `json:"-" db:"host_id"`
	Addr     NullString `json:"addr" db:"addr"`
	AddrType NullString `json:"addrtype" db:"addrtype"`
}

// Port is an internal type that contains a subset of the data in a go-nmap.Port
type Port struct {
	HostID   NullInt64  `json:"-" db:"host_id"`
	Protocol NullString `json:"protocol" db:"protocol"`
	PortID   NullInt64  `json:"portid" db:"port_id"`
	Owner    NullString `json:"owner" db:"owner"`
	Service  NullString `json:"service" db:"service"`
}

// Hostname is an internal type that contains a subset of the data in a go-nmap.Hostname
type Hostname struct {
	HostID NullInt64  `json:"-" db:"host_id"`
	Name   NullString `json:"name" db:"name"`
	Type   NullString `json:"type" db:"type"`
}

// ParseXMLData parses data from incoming nmap XML files to structs
// valid file types are XML
func ParseXMLData(fileID string, data []byte) ([]*Host, error) {
	d, err := nmap.Parse(data)
	if err != nil {
		return nil, err
	}

	var hosts []*Host

	for _, h := range d.Hosts {

		host := &Host{
			FileID:    fileID,
			StartTime: time.Time(h.StartTime),
			EndTime:   time.Time(h.EndTime),
			Comment:   h.Comment,
			Status:    h.Status.Reason,
		}
		// Parse hostnames
		for _, hn := range h.Hostnames {
			hostname := &Hostname{
				Name: ToNullString(hn.Name),
				Type: ToNullString(hn.Type),
			}
			host.Hostnames = append(host.Hostnames, hostname)
		}
		// Parse ports
		for _, p := range h.Ports {
			port := &Port{
				Protocol: ToNullString(p.Protocol),
				PortID:   ToNullInt64(p.PortId),
				Owner:    ToNullString(p.Owner.Name),
				Service:  ToNullString(p.Service.Name),
			}
			host.Ports = append(host.Ports, port)
		}
		// Parse addresses
		for _, a := range h.Addresses {
			address := &Address{
				Addr:     ToNullString(a.Addr),
				AddrType: ToNullString(a.AddrType),
			}
			host.Addresses = append(host.Addresses, address)
		}
		hosts = append(hosts, host)
	}
	return hosts, nil
}

type NullString sql.NullString
type NullInt64 sql.NullInt64

func ToNullString(s string) NullString {
	return NullString{String: s, Valid: s != ""}
}

func (n *NullString) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &n.String)
	n.Valid = (err == nil)
	return err
}

func (n *NullString) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.String)
}

func (n *NullString) Scan(value interface{}) error {
	var s sql.NullString
	if err := s.Scan(value); err != nil {
		return err
	}

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*n = NullString{s.String, false}
	} else {
		*n = NullString{s.String, true}
	}

	return nil
}

func ToNullInt64(i int) NullInt64 {
	return NullInt64{Int64: int64(i), Valid: true}
}

func (n *NullInt64) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &n.Int64)
	n.Valid = (err == nil)
	return err
}

func (n *NullInt64) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.Int64)
}

// Scan implements the Scanner interface for NullInt64
func (n *NullInt64) Scan(value interface{}) error {
	var i sql.NullInt64
	if err := i.Scan(value); err != nil {
		return err
	}

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*n = NullInt64{i.Int64, false}
	} else {
		*n = NullInt64{i.Int64, true}
	}
	return nil
}
