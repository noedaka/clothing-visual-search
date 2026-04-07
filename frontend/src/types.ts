export interface Category {
  id: number;
  name: string;
}

export interface Product {
  id: number;
  name: string;
  description: string;
  price: number;
  category_id: number;
}

export interface Image {
  id: number;
  product_id: number;
  object_key: string;
  is_primary: boolean;
  url: string;
}

export interface ProductWithImages {
  product: Product;
  images: Image[];
}
