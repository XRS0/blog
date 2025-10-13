import { useEffect, useRef, useState } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';
import { deleteArticle, fetchArticle, likeArticle } from '../api/articles';
import type { Article } from '../types';
import { useAuth } from '../context/AuthContext';
import { MarkdownRenderer } from '../components/MarkdownRenderer';

function formatDateTime(dateString: string): string {
  return new Intl.DateTimeFormat('ru-RU', {
    day: '2-digit',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  }).format(new Date(dateString));
}

function ArticleDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [article, setArticle] = useState<Article | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { token, user } = useAuth();
  const lastRequestKey = useRef<string | null>(null);
  const mountedRef = useRef(false);

  useEffect(() => {
    mountedRef.current = true;
    return () => {
      mountedRef.current = false;
    };
  }, []);

  useEffect(() => {
    if (!id) {
      return;
    }

    const key = `${id}:${token ?? ''}`;
    if (lastRequestKey.current === key) {
      return;
    }

    lastRequestKey.current = key;

    (async () => {
      if (!mountedRef.current) {
        return;
      }
      setLoading(true);
      try {
        const data = await fetchArticle(Number(id), token || undefined);
        if (!mountedRef.current) {
          return;
        }
        setArticle(data);
        setError(null);
      } catch (err) {
        if (!mountedRef.current) {
          return;
        }
        setError(err instanceof Error ? err.message : 'Не удалось загрузить статью');
      } finally {
        if (mountedRef.current) {
          setLoading(false);
        }
      }
    })();
  }, [id, token]);

  async function handleLike() {
  if (!article || !token) {
      navigate(`/login?redirect=/article/${id}`);
      return;
    }
    try {
  const { article: updated } = await likeArticle(article.id, token, !article.viewerLiked);
      setArticle(updated);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Не удалось поставить лайк');
    }
  }

  async function handleDelete() {
  if (!article || !token) {
      navigate(`/login?redirect=/article/${id}`);
      return;
    }
    const confirmed = window.confirm('Удалить статью без возможности восстановления?');
    if (!confirmed) {
      return;
    }

    try {
  await deleteArticle(article.id, token);
      navigate('/');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Не удалось удалить статью');
    }
  }

  if (loading) {
    return <p className="helper-text">Загружаем статью…</p>;
  }

  if (error) {
    return (
      <div className="empty-state">
        <p>{error}</p>
        <Link to="/" className="primary-button" style={{ boxShadow: 'none' }}>
          Вернуться к списку
        </Link>
      </div>
    );
  }

  if (!article) {
    return null;
  }

  return (
    <article className="floating-panel" style={{ display: 'grid', gap: '1.8rem', maxWidth: '1200px', margin: '0 auto' }}>
      <div style={{ display: 'grid', gap: '1.1rem' }}>
        <h1 className="page-title">{article.title}</h1>
        <div className="status-bar">
          <span className="status-dot" aria-hidden />
          <span>{article.views} просмотров</span>
          <span>{article.likes} ❤</span>
          <span>Обновлено {formatDateTime(article.updatedAt)}</span>
          {article.author && <span>Автор: {article.author.username}</span>}
        </div>
      </div>
      <div className="article-body markdown-preview">
        <MarkdownRenderer content={article.content} />
      </div>
      <div className="article-actions">
        <button type="button" className="primary-button" onClick={handleLike}>
          {article.viewerLiked ? '❤ Больше не нравится' : '❤ Нравится'}
        </button>
        {user?.id === article.userId && (
          <>
            <Link to={`/article/${article.id}/edit`} className="secondary-button">
              Редактировать
            </Link>
            <button
              type="button"
              onClick={handleDelete}
              className="secondary-button"
              style={{ borderColor: 'rgba(59, 130, 246, 0.4)', color: 'var(--color-accent)' }}
            >
              Удалить
            </button>
          </>
        )}
      </div>
      {!token && <p className="helper-text">Войдите, чтобы ставить лайки и управлять статьями.</p>}
    </article>
  );
}

export default ArticleDetailPage;
