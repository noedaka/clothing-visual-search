import { Link, Outlet, useLocation } from 'react-router-dom';
import './Layout.css';

export function Layout() {
  const location = useLocation();

  const isActive = (path: string) => location.pathname === path;

  return (
    <div className="layout">
      <nav className="navbar">
        <div className="nav-brand">Clothing Visual Search</div>
        <ul className="nav-links">
          <li>
            <Link to="/" className={isActive('/') ? 'active' : ''}>
              Search
            </Link>
          </li>
          <li>
            <Link to="/categories" className={isActive('/categories') ? 'active' : ''}>
              Categories
            </Link>
          </li>
          <li>
            <Link to="/products" className={isActive('/products') ? 'active' : ''}>
              Products
            </Link>
          </li>
        </ul>
      </nav>
      <main className="main-content">
        <Outlet />
      </main>
    </div>
  );
}
