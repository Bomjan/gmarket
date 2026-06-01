package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Bomjan/gmarket/backend/internal/config"
	"github.com/Bomjan/gmarket/backend/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	slog.Info("creating new database")
	db, err := sql.Open("sqlite3", cfg.StoragePath)

	if err != nil {
		return nil, err
	}

	// Create student table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS student (
		id INTEGER PRIMARY KEY AUTOINCREMENT
		,name TEXT
		,email TEXT
		,age INTEGER
	);`)

	if err != nil {
		return nil, err
	}

	// Create product table
	slog.Info("creating product table")
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS product (
	    id INTEGER PRIMARY KEY AUTOINCREMENT
		,name TEXT
	    ,price DECIMAL
	);`)
	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil

}
func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO student(name, email, age) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastId, nil
}
func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT * FROM student WHERE id=? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}

	defer stmt.Close()

	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return types.Student{}, fmt.Errorf("no students found with id %s", fmt.Sprint(id))
		}
		return types.Student{}, fmt.Errorf(err.Error())
	}

	return student, nil

}

func (s *Sqlite) GetAllStudents() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM student")
	if err != nil {
		return []types.Student{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return []types.Student{}, err
	}
	defer rows.Close()

	var students []types.Student
	for rows.Next() {
		var student types.Student
		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return []types.Student{}, err
		}

		students = append(students, student)
	}

	return students, nil
}

func (s *Sqlite) CreateProduct(name string, price float64) (int64, error) {

	slog.Info("creating new product")
	stmt, err := s.Db.Prepare("INSERT INTO product(name, price) VALUES (?, ?)")
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	id, err := stmt.Exec(name, price)
	if err != nil {
		return 0, err
	}
	lastId, err := id.LastInsertId()
	if err != nil {
		return 0, err
	}
	slog.Info("product created", slog.String("id", fmt.Sprint(lastId)))

	return lastId, nil
}
