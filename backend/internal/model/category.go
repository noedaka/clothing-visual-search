package model

type Category struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}
