import { useState, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { searchByImage } from '../api';
import type { ProductWithImages } from '../types';
import './Search.css';

export function Search() {
  const [selectedImage, setSelectedImage] = useState<File | null>(null);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const [searching, setSearching] = useState(false);
  const [results, setResults] = useState<ProductWithImages[]>([]);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  const handleImageChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0];
      setSelectedImage(file);
      setPreviewUrl(URL.createObjectURL(file));
      setResults([]);
      setError(null);
    }
  }, []);

  const clearImage = useCallback(() => {
    if (previewUrl) {
      URL.revokeObjectURL(previewUrl);
    }
    setSelectedImage(null);
    setPreviewUrl(null);
    setResults([]);
    setError(null);
  }, [previewUrl]);

  async function handleSearch(e: React.FormEvent) {
    e.preventDefault();
    if (!selectedImage) return;

    try {
      setSearching(true);
      setError(null);
      const data = await searchByImage(selectedImage);
      setResults(data);
    } catch {
      setError('Search failed. Please try again.');
    } finally {
      setSearching(false);
    }
  }

  function viewProductDetails(productWithImages: ProductWithImages) {
    navigate('/product', { state: { product: productWithImages } });
  }

  return (
    <div className="search-page">
      <h1>Visual Search</h1>
      <p className="subtitle">Upload an image to find similar clothing items</p>

      <form onSubmit={handleSearch} className="search-form">
        <div className="upload-section">
          {!previewUrl ? (
            <label className="upload-area">
              <input
                type="file"
                accept="image/*"
                onChange={handleImageChange}
                className="file-input"
              />
              <div className="upload-content">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="48"
                  height="48"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                >
                  <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
                  <polyline points="17 8 12 3 7 8" />
                  <line x1="12" y1="3" x2="12" y2="15" />
                </svg>
                <p>Click to upload or drag and drop</p>
                <small>Supports: JPG, PNG, WebP</small>
              </div>
            </label>
          ) : (
            <div className="preview-section">
              <img src={previewUrl} alt="Search preview" />
              <button type="button" className="clear-btn" onClick={clearImage}>
                Remove image
              </button>
            </div>
          )}
        </div>

        {previewUrl && (
          <button
            type="submit"
            disabled={searching}
            className="search-btn"
          >
            {searching ? 'Searching...' : 'Search Similar Items'}
          </button>
        )}
      </form>

      {error && <div className="alert error">{error}</div>}

      {results.length > 0 && (
        <div className="results-section">
          <h2>Search Results</h2>
          <p className="results-count">Found {results.length} similar items</p>

          <div className="results-grid">
            {results.map((item, index) => (
              <div
                key={item.product.id}
                className="result-card"
                onClick={() => viewProductDetails(item)}
              >
                <div className="result-image">
                  {item.images.length > 0 ? (
                    <img src={item.images[0].url} alt={item.product.name} />
                  ) : (
                    <div className="no-image">No image</div>
                  )}
                  <span className="match-rank">#{index + 1}</span>
                </div>
                <div className="result-info">
                  <h3>{item.product.name}</h3>
                  <p className="price">${item.product.price.toFixed(2)}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {results.length === 0 && !searching && selectedImage && (
        <div className="no-results">
          <p>No similar items found. Try a different image.</p>
        </div>
      )}
    </div>
  );
}
