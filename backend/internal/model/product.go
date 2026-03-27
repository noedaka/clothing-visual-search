package model

type Product struct {
	ID          int64   `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	Description string  `json:"description" db:"description"`
	Price       float32 `json:"price" db:"price"`
	CategoryID  int64   `json:"category_id" db:"category_id"`
}

type ProductWithImagesData struct {
	Product           Product     `json:"product"`
	ProductImagesData []ImageData `json:"images_data"`
}

type ProductWithImages struct {
	Product       Product `json:"product"`
	ProductImages []Image `json:"images"`
}

type AddproductRequest struct {
	
}