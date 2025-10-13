import { ChangeEvent, FormEvent, useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { createArticle, fetchArticlePreview, updateArticle } from '../api/articles';
import type { ArticlePayload } from '../types';
import { useAuth } from '../context/AuthContext';
import { MarkdownRenderer } from '../components/MarkdownRenderer';

type EditorMode = 'create' | 'edit';
type ViewMode = 'edit' | 'preview';

type ArticleEditorPageProps = {
  mode: EditorMode;
};

const emptyForm: ArticlePayload = {
  title: '',
  content: ''
};

function ArticleEditorPage({ mode }: ArticleEditorPageProps) {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [form, setForm] = useState<ArticlePayload>(emptyForm);
  const [loading, setLoading] = useState(mode === 'edit');
  const [error, setError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);
  const [viewMode, setViewMode] = useState<ViewMode>('edit');
  const { token } = useAuth();

  useEffect(() => {
    if (!token) {
      navigate(
        `/login?redirect=${encodeURIComponent(mode === 'edit' && id ? `/article/${id}/edit` : '/article/new')}`
      );
      return;
    }

    async function loadArticle(articleId: number) {
      try {
        setLoading(true);
  const preview = await fetchArticlePreview(articleId, token!);
        setForm({ title: preview.title, content: preview.content });
        setError(null);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Не удалось загрузить статью для редактирования');
      } finally {
        setLoading(false);
      }
    }

    if (mode === 'edit' && id) {
      loadArticle(Number(id));
    }
  }, [id, mode, navigate, token]);

  function handleChange(event: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) {
    const { name, value } = event.target;
    setForm((prev: ArticlePayload) => ({ ...prev, [name]: value }));
    // Очищаем ошибку при изменении полей
    if (error) {
      setError(null);
    }
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (saving) {
      return;
    }

    const trimmedTitle = form.title.trim();
    const trimmedContent = form.content.trim();

    // Валидация
    if (!trimmedTitle || !trimmedContent) {
      setError('Пожалуйста, заполните название и текст статьи.');
      return;
    }

    if (trimmedTitle.length < 3) {
      setError('Заголовок должен содержать минимум 3 символа.');
      return;
    }

    if (trimmedContent.length < 10) {
      setError('Текст статьи должен содержать минимум 10 символов.');
      return;
    }

    if (!token) {
      navigate(
        `/login?redirect=${encodeURIComponent(mode === 'edit' && id ? `/article/${id}/edit` : '/article/new')}`
      );
      return;
    }

    try {
      setSaving(true);
      setError(null);
      if (mode === 'create') {
        const article = await createArticle(
          {
            title: trimmedTitle,
            content: trimmedContent
          },
          token!
        );
        navigate(`/article/${article.id}`);
      } else if (id) {
        const article = await updateArticle(
          Number(id),
          {
            title: trimmedTitle,
            content: trimmedContent
          },
          token!
        );
        navigate(`/article/${article.id}`);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Не удалось сохранить статью');
    } finally {
      setSaving(false);
    }
  }

  const title = mode === 'create' ? 'Новая статья' : 'Редактирование статьи';

  if (loading) {
    return <p className="helper-text">Подготавливаем редактор…</p>;
  }

  return (
    <section className="floating-panel" style={{ display: 'grid', gap: '1.75rem' }}>
      <div style={{ display: 'grid', gap: '0.65rem' }}>
        <h1 className="page-title">{title}</h1>
        <p className="helper-text" style={{ fontSize: '1rem' }}>
          {mode === 'create'
            ? 'Расскажите историю, поделитесь идеей или оформите гайд — мы подсветим лучшие мысли.'
            : 'Освежите свою публикацию: уточните детали, добавьте новые факты и сохраните изменения.'}
        </p>
      </div>

      <form className="input-stack glass" onSubmit={handleSubmit}>
        <div>
          <label className="label" htmlFor="title">
            Заголовок
          </label>
          <input
            id="title"
            name="title"
            className="input"
            value={form.title}
            onChange={handleChange}
            placeholder="Заголовок, который захочется открыть"
            minLength={3}
            maxLength={120}
            required
          />
          <p className="helper-text">От 3 до 120 символов — будьте доходчивы и интригующие.</p>
        </div>
        <div>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '0.5rem' }}>
            <label className="label" htmlFor="content" style={{ marginBottom: 0 }}>
              Текст статьи
            </label>
            <div className="editor-toggle">
              <button
                type="button"
                className={`editor-toggle-btn ${viewMode === 'edit' ? 'active' : ''}`}
                onClick={() => setViewMode('edit')}
              >
                ✏️ Редактор
              </button>
              <button
                type="button"
                className={`editor-toggle-btn ${viewMode === 'preview' ? 'active' : ''}`}
                onClick={() => setViewMode('preview')}
              >
                👁️ Просмотр
              </button>
            </div>
          </div>
          {viewMode === 'edit' ? (
            <>
              <textarea
                id="content"
                name="content"
                className="textarea"
                value={form.content}
                onChange={handleChange}
                placeholder="Поделитесь вашими мыслями, наблюдениями и историями…

Поддерживается Markdown:
# Заголовок 1
## Заголовок 2
**жирный текст**
*курсив*
- список
1. нумерованный список
[ссылка](https://example.com)
`код`
```
блок кода
```
> цитата
| таблица | данные |
|---------|--------|
| ячейка  | ячейка |"
                minLength={10}
                required
                style={{ minHeight: '400px' }}
              />
              <p className="helper-text">
                Минимум 10 символов. Поддерживается Markdown форматирование.
                {form.content.trim().length > 0 && form.content.trim().length < 10 && (
                  <span style={{ color: 'var(--color-accent)', marginLeft: '0.5rem' }}>
                    (осталось {10 - form.content.trim().length} симв.)
                  </span>
                )}
              </p>
            </>
          ) : (
            <div className="article-body markdown-preview" style={{ minHeight: '400px' }}>
              {form.content.trim() ? (
                <MarkdownRenderer content={form.content} />
              ) : (
                <p className="helper-text" style={{ textAlign: 'center', padding: '3rem' }}>
                  Начните писать, чтобы увидеть предварительный просмотр
                </p>
              )}
            </div>
          )}
        </div>
        {error && (
          <p className="helper-text" style={{ color: 'var(--color-accent)', margin: 0 }} aria-live="assertive">
            {error}
          </p>
        )}
        <div className="article-actions" style={{ marginTop: '0.5rem' }}>
          <button type="submit" className="primary-button" disabled={saving}>
            {saving ? 'Сохраняем…' : mode === 'create' ? 'Опубликовать' : 'Обновить статью'}
          </button>
          <button
            type="button"
            className="secondary-button"
            onClick={() => navigate(mode === 'create' ? '/' : `/article/${id}`)}
            disabled={saving}
          >
            Отменить
          </button>
        </div>
      </form>
    </section>
  );
}

export default ArticleEditorPage;
