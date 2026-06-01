package types

type Student struct {
	Id    int64
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required"`
	Age   int    `json:"age" validate:"required"`
}

type Product struct {
	Id    int64   `json:"id" validate:"required"`
	Name  string  `json:"name" validate:"required"`
	Price float64 `json:"price" validate:"required"`
}
