package storage

import "github.com/Bomjan/gmarket/backend/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetAllStudents() ([]types.Student, error)

	CreateProduct(name string, price float64) (int64, error)
	GetProductById(id int64) (types.Product, error)
	GetAllProducts() ([]types.Product, error)
	DeleteProductById(id int64) error
}
