package model

type Category struct {
	ID   int64  `json:"-" db:"id"`
	Name string `json:"name" db:"name"`
}
