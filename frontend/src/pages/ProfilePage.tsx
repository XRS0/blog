import { FormEvent, useEffect, useMemo, useState } from 'react';
import { useAuth } from '../context/AuthContext';

function ProfilePage() {
  const { user, loading, refreshProfile, logout, updateProfile, error } = useAuth();
  const [form, setForm] = useState({ username: '', contacts: '' });
  const [saving, setSaving] = useState(false);
  const [message, setMessage] = useState<string | null>(null);

  useEffect(() => {
    refreshProfile();
  }, [refreshProfile]);

  useEffect(() => {
    if (user) {
      setForm({
        username: user.username,
        contacts: (Array.isArray(user.contacts) ? user.contacts : []).join('\n')
      });
    }
  }, [user]);

  const parsedContacts = useMemo(
    () =>
      form.contacts
        .split(/\r?\n/)
        .map((value) => value.trim())
        .filter(Boolean),
    [form.contacts]
  );

  const feedbackMessage = message ?? (error ?? null);
  const feedbackColor = message ? '#16a34a' : '#ef4444';

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!user) {
      return;
    }

    try {
      setSaving(true);
      await updateProfile({ username: form.username.trim(), contacts: parsedContacts });
      setMessage('Профиль обновлён');
    } catch (err) {
      setMessage(err instanceof Error ? err.message : 'Не удалось обновить профиль');
    } finally {
      setSaving(false);
    }
  }

  if (loading && !user) {
    return <p className="helper-text">Загружаем профиль…</p>;
  }

  if (!user) {
    return <p className="helper-text">Похоже, вы не авторизованы.</p>;
  }

  return (
    <section className="floating-panel" style={{ maxWidth: '42rem', margin: '2.5rem auto', gap: '1.75rem' }}>
      <header style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', gap: '1rem' }}>
        <div style={{ display: 'grid', gap: '0.35rem' }}>
          <h1 className="page-title" style={{ marginBottom: 0 }}>
            Профиль автора
          </h1>
          <p className="helper-text" style={{ margin: 0 }}>
            Управляйте публичным именем и контактами, чтобы читатели легко нашли вас.
          </p>
        </div>
        <button type="button" className="secondary-button" onClick={logout}>
          Выйти
        </button>
      </header>

      <div className="card glass" style={{ padding: '1.75rem 2rem', gap: '1.35rem' }}>
        <div className="chip-stack">
          <span className="chip">На платформе с {new Date(user.createdAt).toLocaleDateString('ru-RU')}</span>
          <span className="chip">Последнее обновление {new Date(user.updatedAt).toLocaleDateString('ru-RU')}</span>
        </div>
        <form className="input-stack" onSubmit={handleSubmit} style={{ boxShadow: 'none' }}>
          <div>
            <p className="label" style={{ marginBottom: '0.25rem' }}>
              Email
            </p>
            <p className="helper-text" style={{ margin: 0, fontSize: '0.95rem' }}>{user.email}</p>
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
              onChange={(event) => {
                setMessage(null);
                setForm((prev) => ({ ...prev, username: event.target.value }));
              }}
              minLength={3}
              maxLength={32}
              required
            />
            <p className="helper-text">От 3 до 32 символов · покажем его в статьях и профиле.</p>
          </div>
          <div>
            <label className="label" htmlFor="contacts">
              Ссылки на контакты (каждый пункт с новой строки)
            </label>
            <textarea
              id="contacts"
              name="contacts"
              className="textarea"
              value={form.contacts}
              onChange={(event) => {
                setMessage(null);
                setForm((prev) => ({ ...prev, contacts: event.target.value }));
              }}
              placeholder="https://t.me/username"
            />
            <p className="helper-text">
              {parsedContacts.length > 0
                ? `Сохраним ${parsedContacts.length} ${parsedContacts.length === 1 ? 'ссылку' : 'ссылки'} для связи`
                : 'Добавьте ссылки на соцсети или мессенджеры, чтобы читатели могли связаться с вами.'}
            </p>
          </div>
          {feedbackMessage && (
            <p className="helper-text" style={{ color: feedbackColor, margin: 0 }} aria-live="polite">
              {feedbackMessage}
            </p>
          )}
          <div className="article-actions" style={{ marginTop: '0.5rem' }}>
            <button type="submit" className="primary-button" disabled={saving}>
              {saving ? 'Сохраняем…' : 'Сохранить изменения'}
            </button>
          </div>
        </form>
      </div>

      {parsedContacts.length > 0 && (
        <div className="floating-panel glass" style={{ padding: '1.6rem', gap: '1rem' }}>
          <h2 className="card-title" style={{ fontSize: '1.25rem' }}>Контакты</h2>
          <ul style={{ paddingLeft: '1.25rem', margin: 0, display: 'grid', gap: '0.35rem' }}>
            {parsedContacts.map((link) => (
              <li key={link}>
                <a href={link} target="_blank" rel="noopener noreferrer">
                  {link}
                </a>
              </li>
            ))}
          </ul>
        </div>
      )}
    </section>
  );
}

export default ProfilePage;
