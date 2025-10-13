import { Link, NavLink, Route, Routes, useNavigate } from 'react-router-dom';
import ArticleDetailPage from './pages/ArticleDetailPage';
import ArticleEditorPage from './pages/ArticleEditorPage';
import ArticleListPage from './pages/ArticleListPage';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import ProfilePage from './pages/ProfilePage';
import { useAuth } from './context/AuthContext';
import { ThemeToggle } from './components/ThemeToggle';

const navLinkClass = ({ isActive }: { isActive: boolean }) =>
  isActive ? 'nav-link active' : 'nav-link';

function App() {
  const navigate = useNavigate();
  const { user, logout } = useAuth();

  const handleLogout = () => {
    logout();
    navigate('/');
  };

  return (
    <div className="app-shell">
      <header className="header">
        <div className="header-inner">
          <NavLink to="/" className="brand">
            Pulse Notes
          </NavLink>
          <nav className="header-nav">
            <NavLink to="/" className={navLinkClass} end>
              Лента
            </NavLink>
            {user && (
              <NavLink to="/article/new" className={navLinkClass}>
                Новая запись
              </NavLink>
            )}
            {user && (
              <NavLink to="/profile" className={navLinkClass}>
                Профиль
              </NavLink>
            )}
          </nav>
          <div className="header-actions">
            <ThemeToggle />
            {user ? (
              <>
                <span className="helper-text user-email" style={{ margin: 0 }}>{user.email}</span>
                <button type="button" className="secondary-button" onClick={handleLogout}>
                  Выйти
                </button>
              </>
            ) : (
              <>
                <NavLink to="/login" className={navLinkClass}>
                  Войти
                </NavLink>
                <Link to="/register" className="primary-button">
                  Регистрация
                </Link>
              </>
            )}
          </div>
        </div>
      </header>

      <main className="container">
        <Routes>
          <Route path="/" element={<ArticleListPage />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/profile" element={<ProfilePage />} />
          <Route path="/article/new" element={<ArticleEditorPage mode="create" />} />
          <Route path="/article/:id" element={<ArticleDetailPage />} />
          <Route path="/article/:id/edit" element={<ArticleEditorPage mode="edit" />} />
          <Route path="*" element={<ArticleListPage />} />
        </Routes>
      </main>
    </div>
  );
}

export default App;
