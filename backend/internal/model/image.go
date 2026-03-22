package model

import "io"

type Image struct {
	ID        int64  `db:"id"`
	ProductID int64  `db:"product_id"`
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

	Embedding []float32
}
