package sqlite

import (
	"database/sql"

	"github.com/Bomjan/gmarket/backend/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)

	if err != nil {
		return nil, err
	}

	db.Exec(`CREATE TABLE IF NOT EXISTS student (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT,
	age INTEGER
	);`)

	return &Sqlite{
		Db: db,
	}, nil

}
func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	s.Db.Prepare("INSERT INTO students(name, email, age) VALUES (?, ?, ?)")
	return 0, nil
}
