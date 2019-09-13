package sqlite

import "database/sql"

type HostRepo struct {
	DB *sql.DB
}
