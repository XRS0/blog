import { FormEvent, useState } from 'react';
import { Link, useNavigate, useSearchParams } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

function LoginPage() {
  const { login, loading, error } = useAuth();
  const [form, setForm] = useState({ email: '', password: '' });
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    try {
      await login(form);
      const redirect = searchParams.get('redirect') || '/';
      navigate(redirect);
    } catch (err) {
      // error is handled in context state
    }
  }

  return (
    <section className="floating-panel" style={{ maxWidth: '28rem', margin: '2.5rem auto', gap: '1.75rem' }}>
      <div style={{ display: 'grid', gap: '0.6rem' }}>
        <h1 className="page-title">С возвращением!</h1>
        <p className="helper-text" style={{ fontSize: '1rem' }}>
          Войдите в аккаунт, чтобы продолжить писать, комментировать и следить за любимыми авторами.
        </p>
      </div>

      <form className="input-stack glass" onSubmit={handleSubmit}>
        <div>
          <label className="label" htmlFor="email">
            Email
          </label>
          <input
            id="email"
            name="email"
            className="input"
            type="email"
            value={form.email}
            onChange={(event) => setForm((prev) => ({ ...prev, email: event.target.value }))}
            placeholder="name@example.com"
            required
            autoComplete="email"
          />
        </div>
        <div>
          <label className="label" htmlFor="password">
            Пароль
          </label>
          <input
            id="password"
            name="password"
            className="input"
            type="password"
            value={form.password}
            onChange={(event) => setForm((prev) => ({ ...prev, password: event.target.value }))}
            placeholder="••••••••"
            required
            autoComplete="current-password"
          />
        </div>
        {error && (
          <p className="helper-text" style={{ color: 'var(--color-accent)', margin: 0 }} aria-live="assertive">
            {error}
          </p>
        )}
        <button type="submit" className="primary-button" disabled={loading}>
          {loading ? 'Входим…' : 'Войти'}
        </button>
        <p className="helper-text" style={{ textAlign: 'center', marginBottom: 0 }}>
          Нет аккаунта? <Link to="/register">Зарегистрируйтесь</Link>
        </p>
      </form>
    </section>
  );
}

export default LoginPage;
