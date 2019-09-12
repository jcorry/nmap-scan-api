package models

import "time"

// Host is an internal type that contains a subset of the data in a go-nmap.Host
type Host struct {
	FileID    string     `json:"fileid" db:"file_id"`
	StartTime time.Time  `json:"starttime" db:"starttime"`
	EndTime   time.Time  `json:"endtime" db:"endtime"`
	Comment   string     `json:"comment" db:"comment"`
	Status    string     `json:"status" db:"status"`
	Hostnames []Hostname `json:"hostnames"`
	Addresses []Address  `json:"addresses"`
	Ports     []Port     `json:"ports"`
}

// Address is an internal type that contains a subset of the data in a go-nmap.Address
type Address struct {
	HostID   int    `json:"hostid" db:"host_id"`
	Addr     string `json:"addr" db:"addr"`
	AddrType string `json:"addrtype" db:"addrtype"`
}

// Port is an internal type that contains a subset of the data in a go-nmap.Port
type Port struct {
	HostID   int    `json:"hostid" db:"host_id"`
	Protocol string `json:"protocol" db:"protocol"`
	PortID   int    `json:"portid" db:"port_id"`
	Owner    string `json:"owner" db:"owner"`
	Service  string `json:"service" db:"service"`
}

// Hostname is an internal type that contains a subset of the data in a go-nmap.Hostname
type Hostname struct {
	HostID int    `json:"hostid" db:"host_id"`
	Name   string `json:"name" db:"name"`
	Type   string `json:"type" db:"type"`
}
