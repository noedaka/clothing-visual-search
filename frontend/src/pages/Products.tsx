import { useState, useEffect, useCallback } from 'react';
import { listCategories, addProduct } from '../api';
import type { Category, Product } from '../types';
import './Products.css';

export function Products() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const [product, setProduct] = useState<Partial<Product>>({
    name: '',
    description: '',
    price: 0,
    category_id: 0,
  });
  const [images, setImages] = useState<File[]>([]);
  const [imagePreviewUrls, setImagePreviewUrls] = useState<string[]>([]);

  useEffect(() => {
    loadCategories();
  }, []);

  async function loadCategories() {
    try {
      const data = await listCategories();
      setCategories(data);
      if (data.length > 0) {
        setProduct((p) => ({ ...p, category_id: data[0].id }));
      }
    } catch {
      setError('Failed to load categories');
    } finally {
      setLoading(false);
    }
  }

  const handleImageChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      const files = Array.from(e.target.files);
      setImages(files);
      const urls = files.map((file) => URL.createObjectURL(file));
      setImagePreviewUrls(urls);
    }
  }, []);

  const removeImage = useCallback((index: number) => {
    setImages((prev) => prev.filter((_, i) => i !== index));
    setImagePreviewUrls((prev) => {
      URL.revokeObjectURL(prev[index]);
      return prev.filter((_, i) => i !== index);
    });
  }, []);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();

    if (!product.name || (product.price || 0) <= 0 || (product.category_id || 0) <= 0) {
      setError('Please fill in all required fields');
      return;
    }

    if (images.length === 0) {
      setError('Please add at least one image');
      return;
    }

    try {
      setSubmitting(true);
      setError(null);

      const productData: Product = {
        id: 0,
        name: product.name,
        description: product.description || '',
        price: product.price || 0,
        category_id: product.category_id || 0,
      };

      await addProduct(productData, images);

      setSuccess('Product added successfully');
      setProduct({
        name: '',
        description: '',
        price: 0,
        category_id: categories[0]?.id || 0,
      });
      setImages([]);
      imagePreviewUrls.forEach(URL.revokeObjectURL);
      setImagePreviewUrls([]);

      setTimeout(() => setSuccess(null), 3000);
    } catch {
      setError('Failed to add product');
    } finally {
      setSubmitting(false);
    }
  }

  if (loading) {
    return <p className="loading">Loading...</p>;
  }

  return (
    <div className="products-page">
      <h1>Add Product</h1>

      {categories.length === 0 ? (
        <div className="alert error">
          Please add at least one category before adding products.
        </div>
      ) : (
        <form onSubmit={handleSubmit} className="product-form">
          <div className="form-row">
            <div className="form-group">
              <label htmlFor="name">
                Product Name <span className="required">*</span>
              </label>
              <input
                type="text"
                id="name"
                value={product.name}
                onChange={(e) => setProduct({ ...product, name: e.target.value })}
                placeholder="Enter product name"
                required
              />
            </div>

            <div className="form-group">
              <label htmlFor="price">
                Price <span className="required">*</span>
              </label>
              <input
                type="number"
                id="price"
                min="0.01"
                step="0.01"
                value={product.price || ''}
                onChange={(e) => setProduct({ ...product, price: parseFloat(e.target.value) || 0 })}
                required
              />
            </div>
          </div>

          <div className="form-group">
            <label htmlFor="category">
              Category <span className="required">*</span>
            </label>
            <select
              id="category"
              value={product.category_id ?? ''}
              onChange={(e) => setProduct({ ...product, category_id: parseInt(e.target.value) })}
              required
            >
              {categories.map((cat) => (
                <option key={cat.id} value={cat.id}>
                  {cat.name}
                </option>
              ))}
            </select>
          </div>

          <div className="form-group">
            <label htmlFor="description">Description</label>
            <textarea
              id="description"
              rows={4}
              value={product.description}
              onChange={(e) => setProduct({ ...product, description: e.target.value })}
              placeholder="Enter product description"
            />
          </div>

          <div className="form-group">
            <label htmlFor="images">
              Images <span className="required">*</span>
            </label>
            <input
              type="file"
              id="images"
              accept="image/*"
              multiple
              onChange={handleImageChange}
              className="file-input"
            />
            <small className="hint">First image will be set as primary</small>
          </div>

          {imagePreviewUrls.length > 0 && (
            <div className="image-previews">
              {imagePreviewUrls.map((url, index) => (
                <div key={index} className={`preview-item ${index === 0 ? 'primary' : ''}`}>
                  <img src={url} alt={`Preview ${index + 1}`} />
                  {index === 0 && <span className="primary-badge">Primary</span>}
                  <button
                    type="button"
                    className="remove-btn"
                    onClick={() => removeImage(index)}
                  >
                    ×
                  </button>
                </div>
              ))}
            </div>
          )}

          {error && <div className="alert error">{error}</div>}
          {success && <div className="alert success">{success}</div>}

          <button type="submit" disabled={submitting} className="submit-btn">
            {submitting ? 'Adding Product...' : 'Add Product'}
          </button>
        </form>
      )}
    </div>
  );
}
