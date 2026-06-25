import { useState, useCallback, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { searchByImage } from '../api';
import type { ProductWithImages } from '../types';
import './Search.css';

const SEARCH_STATE_KEY = 'visualSearchState';

export function Search() {
  const [selectedImage, setSelectedImage] = useState<File | null>(null);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const [previewBase64, setPreviewBase64] = useState<string | null>(null);
  const [searching, setSearching] = useState(false);
  const [results, setResults] = useState<ProductWithImages[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [hasSearched, setHasSearched] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    const saved = sessionStorage.getItem(SEARCH_STATE_KEY);
    if (saved) {
      try {
        const state = JSON.parse(saved);
        setResults(state.results || []);
        setHasSearched(state.hasSearched || false);
        if (state.previewBase64) {
          setPreviewUrl(state.previewBase64);
          setPreviewBase64(state.previewBase64);
        }
      } catch (e) {
        console.error('Failed to restore search state', e);
      }
    }
  }, []);

  useEffect(() => {
    if (results.length > 0 || previewBase64) {
      const state = {
        results,
        previewBase64,
        hasSearched,
      };
      sessionStorage.setItem(SEARCH_STATE_KEY, JSON.stringify(state));
    }
  }, [results, previewBase64, hasSearched]);

  const handleImageChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0];
      setSelectedImage(file);
      const blobUrl = URL.createObjectURL(file);
      setPreviewUrl(blobUrl);
      setResults([]);
      setError(null);
      setHasSearched(false);

      const reader = new FileReader();
      reader.onloadend = () => {
        setPreviewBase64(reader.result as string);
      };
      reader.readAsDataURL(file);
    }
  }, []);

  const clearImage = useCallback(() => {
    if (previewUrl && previewUrl.startsWith('blob:')) {
      URL.revokeObjectURL(previewUrl);
    }
    setSelectedImage(null);
    setPreviewUrl(null);
    setPreviewBase64(null);
    setResults([]);
    setError(null);
    setHasSearched(false);
    sessionStorage.removeItem(SEARCH_STATE_KEY);
  }, [previewUrl]);

  async function handleSearch(e: React.FormEvent) {
    e.preventDefault();
    if (!selectedImage) return;

    try {
      setSearching(true);
      setError(null);
      const data = await searchByImage(selectedImage);
      setResults(data);
      setHasSearched(true);
    } catch {
      setError('Поиск не удался. Попробуйте ещё раз.');
      setHasSearched(true);
    } finally {
      setSearching(false);
    }
  }

  function viewProductDetails(productWithImages: ProductWithImages) {
    navigate('/product', { state: { product: productWithImages } });
  }

  return (
    <div className="search-page">
      <h1>Визуальный поиск</h1>
      <p className="subtitle">Загрузите фотографию, чтобы найти похожие товары</p>

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
                <p>Нажмите для загрузки или перетащите файл</p>
                <small>Форматы: JPG, PNG, WebP</small>
              </div>
            </label>
          ) : (
            <div className="preview-section">
              <img src={previewUrl} alt="Search preview" />
              <button type="button" className="clear-btn" onClick={clearImage}>
                Удалить фото
              </button>
            </div>
          )}
        </div>

        {selectedImage && previewUrl && (
          <button type="submit" disabled={searching} className="search-btn">
            {searching ? 'Поиск...' : 'Найти похожие товары'}
          </button>
        )}
      </form>

      {error && <div className="alert error">{error}</div>}

      {results.length > 0 && (
        <div className="results-section">
          <h2>Результаты поиска</h2>
          <p className="results-count">Найдено {results.length} похожих товаров</p>
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
                    <div className="no-image">Фото отсутствует</div>
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

      {results.length === 0 && !searching && hasSearched && (
        <div className="no-results">
          <p>Похожих товаров не найдено. Попробуйте другое фото.</p>
        </div>
      )}
    </div>
  );
}