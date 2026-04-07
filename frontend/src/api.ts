import type { Category, Product, ProductWithImages } from './types';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

export async function listCategories(): Promise<Category[]> {
  const response = await fetch(`${API_BASE_URL}/category`);
  if (response.status === 204) {
    return [];
  }
  if (!response.ok) {
    throw new Error('Failed to fetch categories');
  }
  return response.json();
}

export async function addCategory(name: string): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/category`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name }),
  });
  if (!response.ok) {
    throw new Error('Failed to add category');
  }
}

export async function addProduct(product: Product, images: File[]): Promise<void> {
  const formData = new FormData();
  formData.append('product', JSON.stringify(product));
  images.forEach((image) => {
    formData.append('images', image);
  });

  const response = await fetch(`${API_BASE_URL}/product`, {
    method: 'POST',
    body: formData,
  });
  if (!response.ok) {
    throw new Error('Failed to add product');
  }
}

export async function searchByImage(image: File): Promise<ProductWithImages[]> {
  const formData = new FormData();
  formData.append('image', image);

  const response = await fetch(`${API_BASE_URL}/product/search`, {
    method: 'POST',
    body: formData,
  });
  if (response.status === 204) {
    return [];
  }
  if (!response.ok) {
    throw new Error('Search failed');
  }
  return response.json();
}
