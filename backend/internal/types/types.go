package types

type Student struct {
	Id    int64
	Name  string `validate:"required"`
	Email string `validate:"required"`
	Age   int    `validate:"required"`
}

type Product struct {
	Id    int64
	Name  string  `validate:"required"`
	Price float64 `validate:"required"`
}
