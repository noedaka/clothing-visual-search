package model

import "io"

type Image struct {
	ID        int    `db:"id"`
	ProductID int    `db:"product_id"`
	ObjectKey string `db:"object_key"`
	IsPrimary bool   `db:"is_primary"`
	URL       string `db:"-"`
}

type ImageData struct {
	File        io.Reader
	FileSize    int64
	Filename    string
	ContentType string
	IsPrimary   bool
}
