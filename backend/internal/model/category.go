package model

type Category struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}
