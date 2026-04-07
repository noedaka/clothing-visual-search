import { useLocation, useNavigate } from 'react-router-dom';
import type { ProductWithImages } from '../types';
import './ProductDetails.css';

interface LocationState {
  product: ProductWithImages;
}

export function ProductDetails() {
  const location = useLocation();
  const navigate = useNavigate();
  const state = location.state as LocationState | null;
  const productWithImages = state?.product;

  if (!productWithImages) {
    return (
      <div className="product-details-page">
        <div className="error-state">
          <p>No product selected</p>
          <button onClick={() => navigate('/')} className="back-btn">
            Go to Search
          </button>
        </div>
      </div>
    );
  }

  const { product, images } = productWithImages;
  const primaryImage = images.find((img) => img.is_primary) || images[0];

  return (
    <div className="product-details-page">
      <button onClick={() => navigate(-1)} className="back-link">
        ← Back to results
      </button>

      <div className="product-content">
        <div className="product-gallery">
          {primaryImage ? (
            <img
              src={primaryImage.url}
              alt={product.name}
              className="main-image"
            />
          ) : (
            <div className="no-image-large">No image available</div>
          )}

          {images.length > 1 && (
            <div className="thumbnail-list">
              {images.map((image) => (
                <img
                  key={image.id}
                  src={image.url}
                  alt={`${product.name} thumbnail`}
                  className={`thumbnail ${image.is_primary ? 'active' : ''}`}
                />
              ))}
            </div>
          )}
        </div>

        <div className="product-info">
          <h1>{product.name}</h1>
          <p className="price">${product.price.toFixed(2)}</p>

          {product.description && (
            <div className="description-section">
              <h2>Description</h2>
              <p>{product.description}</p>
            </div>
          )}

          <div className="details-section">
            <h2>Details</h2>
            <dl>
              <dt>Product ID</dt>
              <dd>{product.id}</dd>
              <dt>Category ID</dt>
              <dd>{product.category_id}</dd>
            </dl>
          </div>
        </div>
      </div>
    </div>
  );
}
