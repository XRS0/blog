import { FormEvent, useMemo, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

function RegisterPage() {
  const { register, loading, error } = useAuth();
  const [form, setForm] = useState({ email: '', password: '', username: '', contacts: '' });
  const navigate = useNavigate();

  const parsedContacts = useMemo(
    () =>
      form.contacts
        .split(/\r?\n/)
        .map((value) => value.trim())
        .filter(Boolean),
    [form.contacts]
  );

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    try {
      await register({
        email: form.email.trim(),
        password: form.password,
        username: form.username.trim(),
        contacts: parsedContacts
      });
      navigate('/');
    } catch (err) {
      // error handled by context state
    }
  }

  return (
    <section className="floating-panel" style={{ maxWidth: '32rem', margin: '2.5rem auto', gap: '1.75rem' }}>
      <div style={{ display: 'grid', gap: '0.6rem' }}>
        <h1 className="page-title">Создайте профиль автора</h1>
        <p className="helper-text" style={{ fontSize: '1rem' }}>
          Несколько полей — и у вас появится собственное пространство для публикаций и общения с аудиторией.
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
          <label className="label" htmlFor="username">
            Имя пользователя
          </label>
          <input
            id="username"
            name="username"
            className="input"
            value={form.username}
            onChange={(event) => setForm((prev) => ({ ...prev, username: event.target.value }))}
            placeholder="pulse_writer"
            minLength={3}
            maxLength={32}
            required
            autoComplete="username"
          />
          <p className="helper-text">От 3 до 32 символов · латиница, цифры и подчёркивания.</p>
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
            placeholder="минимум 6 символов"
            required
            autoComplete="new-password"
          />
        </div>
        <div>
          <label className="label" htmlFor="contacts">
            Контакты (каждый с новой строки)
          </label>
          <textarea
            id="contacts"
            name="contacts"
            className="textarea"
            value={form.contacts}
            onChange={(event) => setForm((prev) => ({ ...prev, contacts: event.target.value }))}
            placeholder="https://t.me/username"
          />
          <p className="helper-text" style={{ marginTop: parsedContacts.length > 0 ? '-0.15rem' : '0.35rem' }}>
            {parsedContacts.length > 0
              ? `Мы сохраним ${parsedContacts.length} ${parsedContacts.length === 1 ? 'ссылку' : 'ссылки'}`
              : 'Добавьте социальные сети, чтобы читатели могли связаться с вами.'}
          </p>
        </div>
        {error && (
          <p className="helper-text" style={{ color: 'var(--color-accent)', margin: 0 }} aria-live="assertive">
            {error}
          </p>
        )}
        <button type="submit" className="primary-button" disabled={loading}>
          {loading ? 'Создаём…' : 'Создать аккаунт'}
        </button>
        <p className="helper-text" style={{ textAlign: 'center', marginBottom: 0 }}>
          Уже зарегистрированы? <Link to="/login">Войти</Link>
        </p>
      </form>
    </section>
  );
}

export default RegisterPage;
