package model

import "io"

type Image struct {
	ID        int64  `json:"id" db:"id"`
	ProductID int64  `json:"product_id" db:"product_id"`
	ObjectKey string `json:"object_key" db:"object_key"`
	IsPrimary bool   `json:"is_primary" db:"is_primary"`
	URL       string `json:"url" db:"-"`
}

type ImageData struct {
	File        io.Reader
	FileSize    int64
	Filename    string
	ContentType string
	IsPrimary   bool

	Embedding []float32
}
