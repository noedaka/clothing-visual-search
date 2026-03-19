package model

type Product struct {
	ID          int     `db:"id"`
	Name        string  `db:"name"`
	Description string  `db:"description"`
	Price       float32 `db:"price"`
	CategoryID  int     `db:"category_id"`
}

type ProductWithImages struct {
	Product       Product
	ProductImages []Image
}
