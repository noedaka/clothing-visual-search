package model

type Product struct {
	ID          int     `db:"id"`
	Name        string  `db:"name"`
	Description string  `db:"description"`
	Price       float32 `db:"price"`
	CategoryID  int     `db:"category_id"`
}

// Добавить конструкторы?
type ProductWithImagesData struct {
	Product           Product
	ProductImagesData []ImageData
}

type ProductWithImages struct {
	Product       Product
	ProductImages []Image
}
