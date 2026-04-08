import { useState, useEffect } from 'react';
import { listCategories, addCategory } from '../api';
import type { Category } from '../types';
import './Categories.css';

export function Categories() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [newCategoryName, setNewCategoryName] = useState('');
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  useEffect(() => {
    loadCategories();
  }, []);

  async function loadCategories() {
    try {
      setLoading(true);
      const data = await listCategories();
      setCategories(data);
      setError(null);
    } catch {
      setError('Failed to load categories');
    } finally {
      setLoading(false);
    }
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!newCategoryName.trim()) return;

    try {
      setSubmitting(true);
      await addCategory(newCategoryName.trim());
      setNewCategoryName('');
      setSuccess('Category added successfully');
      await loadCategories();
      setTimeout(() => setSuccess(null), 3000);
    } catch {
      setError('Failed to add category');
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="categories-page">
      <h1>Categories</h1>

      <div className="add-category-section">
        <h2>Add New Category</h2>
        <form onSubmit={handleSubmit} className="category-form">
          <div className="form-group">
            <label htmlFor="categoryName">Category Name</label>
            <input
              type="text"
              id="categoryName"
              value={newCategoryName}
              onChange={(e) => setNewCategoryName(e.target.value)}
              placeholder="Enter category name"
              required
            />
          </div>
          <button type="submit" disabled={submitting || !newCategoryName.trim()}>
            {submitting ? 'Adding...' : 'Add Category'}
          </button>
        </form>
      </div>

      {error && <div className="alert error">{error}</div>}
      {success && <div className="alert success">{success}</div>}

      <div className="categories-list-section">
        <h2>Existing Categories</h2>
        {loading ? (
          <p className="loading">Loading categories...</p>
        ) : !categories || categories.length === 0 ? (
          <p className="empty">No categories found. Add one above!</p>
        ) : (
          <ul className="categories-list">
            {categories.map((category) => (
              <li key={category.id} className="category-item">
                <span className="category-name">{category.name}</span>
                <span className="category-id">ID: {category.id}</span>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}
