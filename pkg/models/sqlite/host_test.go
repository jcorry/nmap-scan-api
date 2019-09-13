package sqlite_test

import (
	"encoding/json"
	"testing"

	"github.com/jcorry/nmap-scan-api/pkg/models"

	"github.com/jcorry/nmap-scan-api/pkg/models/sqlite"
)

func Test_HostBatchInsert(t *testing.T) {
	db, teardown := newTestDB(t)
	defer teardown()

	tests := []struct {
		name    string
		wantErr error
	}{
		{
			name:    "Successful Insert",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := sqlite.HostRepo{DB: db}

			hostData, err := getValidHostSlice(t)
			if err != nil {
				t.Fatal(err)
			}
			err = h.BatchInsert(hostData)

			// No errs
			if err != tt.wantErr {
				t.Fatalf("Want err %s, Got err %s", tt.wantErr, err)
			}
		})
	}

}

func Test_HostInsert(t *testing.T) {

	db, teardown := newTestDB(t)
	defer teardown()

	tests := []struct {
		name    string
		wantErr error
	}{
		{
			name:    "Successful Insert",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := sqlite.HostRepo{DB: db}

			validHost, err := getValidHost(t)
			if err != nil {
				t.Fatal(err)
			}
			// Perform the insert
			err = h.Insert(validHost)

			// No errs
			if err != tt.wantErr {
				t.Fatalf("Want err %s, Got err %s", tt.wantErr, err)
			}

			// Host ID should not be 0
			if validHost.ID == 0 {
				t.Fatalf("ID not set on host")
			}

			// Host ID should match each port host ID
			for _, p := range validHost.Ports {
				if p.HostID != validHost.ID {
					t.Fatalf("Want port.HostID: %d, Got port.HostID: %d", validHost.ID, p.HostID)
				}
			}
		})
	}
}

func getValidHostSlice(t *testing.T) ([]*models.Host, error) {
	hostJSON := `[
    {
        "fileid": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
        "starttime": "2018-04-20T12:36:54-04:00",
        "endtime": "2018-04-20T12:36:56-04:00",
        "comment": "",
        "status": "user-set",
        "hostnames": [
            {
                "hostid": 0,
                "name": "cpc123026-glen5-2-0-cust970.2-1.cable.virginm.net",
                "type": "PTR"
            }
        ],
        "addresses": [
            {
                "hostid": 0,
                "addr": "81.107.115.203",
                "addrtype": "ipv4"
            }
        ],
        "ports": [
            {
                "hostid": 0,
                "protocol": "tcp",
                "portid": 80,
                "owner": "",
                "service": "http"
            },
            {
                "hostid": 0,
                "protocol": "tcp",
                "portid": 443,
                "owner": "",
                "service": "https"
            },
            {
                "hostid": 0,
                "protocol": "tcp",
                "portid": 5000,
                "owner": "",
                "service": "upnp"
            },
            {
                "hostid": 0,
                "protocol": "tcp",
                "portid": 8080,
                "owner": "",
                "service": "http-proxy"
            },
            {
                "hostid": 0,
                "protocol": "tcp",
                "portid": 8443,
                "owner": "",
                "service": "https-alt"
            }
        ]
    },
    {
        "fileid": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
        "starttime": "2018-04-20T12:36:54-04:00",
        "endtime": "2018-04-20T12:36:56-04:00",
        "comment": "",
        "status": "user-set",
        "hostnames": [
            {
                "hostid": 0,
                "name": "nicholas.cybershark.net",
                "type": "PTR"
            }
        ],
        "addresses": [
            {
                "hostid": 0,
                "addr": "158.69.205.102",
                "addrtype": "ipv4"
            }
        ],
        "ports": [
            {
                "hostid": 0,
                "protocol": "tcp",
                "portid": 80,
                "owner": "",
                "service": "http"
            },
            {
                "hostid": 0,
                "protocol": "tcp",
                "portid": 443,
                "owner": "",
                "service": "https"
            },
            {
                "hostid": 0,
                "protocol": "tcp",
                "portid": 5000,
                "owner": "",
                "service": "upnp"
            },
            {
                "hostid": 0,
                "protocol": "tcp",
                "portid": 8080,
                "owner": "",
                "service": "http-proxy"
            },
            {
                "hostid": 0,
                "protocol": "tcp",
                "portid": 8443,
                "owner": "",
                "service": "https-alt"
            }
        ]
    },
    {
        "fileid": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
        "starttime": "2018-04-20T12:36:54-04:00",
        "endtime": "2018-04-20T12:36:56-04:00",
        "comment": "",
        "status": "user-set",
        "hostnames": [
            {
                "hostid": 0,
                "name": "loghermes.sysraildata.com",
                "type": "PTR"
            }
        ],
        "addresses": [
            {
                "hostid": 0,
                "addr": "193.22.92.195",
                "addrtype": "ipv4"
            }
        ],
        "ports": [
            {
                "hostid": 0,
                "protocol": "tcp",
                "portid": 80,
                "owner": "",
                "service": "http"
            },
            {
                "hostid": 0,
                "protocol": "tcp",
                "portid": 443,
                "owner": "",
                "service": "https"
            },
            {
                "hostid": 0,
                "protocol": "tcp",
                "portid": 5000,
                "owner": "",
                "service": "upnp"
            },
            {
                "hostid": 0,
                "protocol": "tcp",
                "portid": 8080,
                "owner": "",
                "service": "http-proxy"
            },
            {
                "hostid": 0,
                "protocol": "tcp",
                "portid": 8443,
                "owner": "",
                "service": "https-alt"
            }
        ]
    }
]`

	var hosts []*models.Host
	if err := json.Unmarshal([]byte(hostJSON), &hosts); err != nil {
		t.Fatal(err)
	}
	return hosts, nil
}

func getValidHost(t *testing.T) (*models.Host, error) {
	hostJSON := `{
        "fileid": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
        "starttime": "2018-04-20T12:36:54-04:00",
        "endtime": "2018-04-20T12:36:56-04:00",
        "comment": "",
        "status": "user-set",
        "hostnames": [
            {
                "hostid": 0,
                "name": "cpc123026-glen5-2-0-cust970.2-1.cable.virginm.net",
                "type": "PTR"
            }
        ],
        "addresses": [
            {
                "hostid": 0,
                "addr": "81.107.115.203",
                "addrtype": "ipv4"
            }
        ],
        "ports": [
            {
                "hostid": 0,
                "protocol": "tcp",
                "portid": 80,
                "owner": "",
                "service": "http"
            },
            {
                "hostid": 0,
                "protocol": "tcp",
                "portid": 443,
                "owner": "",
                "service": "https"
            },
            {
                "hostid": 0,
                "protocol": "tcp",
                "portid": 5000,
                "owner": "",
                "service": "upnp"
            },
            {
                "hostid": 0,
                "protocol": "tcp",
                "portid": 8080,
                "owner": "",
                "service": "http-proxy"
            },
            {
                "hostid": 0,
                "protocol": "tcp",
                "portid": 8443,
                "owner": "",
                "service": "https-alt"
            }
        ]
    }`
	var host models.Host
	if err := json.Unmarshal([]byte(hostJSON), &host); err != nil {
		t.Fatal(err)
	}
	return &host, nil
}
