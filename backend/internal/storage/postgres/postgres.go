package postgres

import (
	"database/sql"
)

type Postgres struct {
	Db *sql.DB
}

//func New(cfg config.Config) (*Postgres, error) {
//	return
//}
