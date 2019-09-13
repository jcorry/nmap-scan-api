package models

import (
	"time"

	nmap "github.com/tomsteele/go-nmap"
)

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
				Name: hn.Name,
				Type: hn.Type,
			}
			host.Hostnames = append(host.Hostnames, *hostname)
		}
		// Parse ports
		for _, p := range h.Ports {
			port := &Port{
				Protocol: p.Protocol,
				PortID:   p.PortId,
				Owner:    p.Owner.Name,
				Service:  p.Service.Name,
			}
			host.Ports = append(host.Ports, *port)
		}
		// Parse addresses
		for _, a := range h.Addresses {
			address := &Address{
				Addr:     a.Addr,
				AddrType: a.AddrType,
			}
			host.Addresses = append(host.Addresses, *address)
		}
		hosts = append(hosts, host)
	}
	return hosts, nil
}
