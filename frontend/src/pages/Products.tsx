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
      setError('Не удалось загрузить категории');
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
      setError('Заполните все обязательные поля');
      return;
    }

    if (images.length === 0) {
      setError('Добавьте хотя бы одно изображение');
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

      setSuccess('Товар успешно добавлен');
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
      setError('Не удалось добавить товар');
    } finally {
      setSubmitting(false);
    }
  }

  if (loading) {
    return <p className="loading">Загрузка...</p>;
  }

  return (
    <div className="products-page">
      <h1>Добавить товар</h1>

      {categories.length === 0 ? (
        <div className="alert error">
          Сначала добавьте хотя бы одну категорию.
        </div>
      ) : (
        <form onSubmit={handleSubmit} className="product-form">
          <div className="form-row">
            <div className="form-group">
              <label htmlFor="name">
                Название <span className="required">*</span>
              </label>
              <input
                type="text"
                id="name"
                value={product.name}
                onChange={(e) => setProduct({ ...product, name: e.target.value })}
                placeholder="Введите название товара"
                required
              />
            </div>

            <div className="form-group">
              <label htmlFor="price">
                Цена <span className="required">*</span>
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
              Категория <span className="required">*</span>
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
            <label htmlFor="description">Описание</label>
            <textarea
              id="description"
              rows={4}
              value={product.description}
              onChange={(e) => setProduct({ ...product, description: e.target.value })}
              placeholder="Введите описание товара"
            />
          </div>

          <div className="form-group">
            <label htmlFor="images">
              Фотографии <span className="required">*</span>
            </label>
            <input
              type="file"
              id="images"
              accept="image/*"
              multiple
              onChange={handleImageChange}
              className="file-input"
            />
            <small className="hint">Первое фото будет главным</small>
          </div>

          {imagePreviewUrls.length > 0 && (
            <div className="image-previews">
              {imagePreviewUrls.map((url, index) => (
                <div key={index} className={`preview-item ${index === 0 ? 'primary' : ''}`}>
                  <img src={url} alt={`Preview ${index + 1}`} />
                  {index === 0 && <span className="primary-badge">Главное</span>}
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
            {submitting ? 'Добавление...' : 'Добавить товар'}
          </button>
        </form>
      )}
    </div>
  );
}
