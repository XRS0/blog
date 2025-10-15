import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { fetchArticles } from '../api/articles';
import type { Article } from '../types';
import { useAuth } from '../context/AuthContext';

function formatDate(dateString: string): string {
  if (!dateString) {
    return 'Неизвестно';
  }
  
  try {
    const date = new Date(dateString);
    // Проверка на валидность даты
    if (isNaN(date.getTime())) {
      return 'Неизвестно';
    }
    
    return new Intl.DateTimeFormat('ru-RU', {
      day: '2-digit',
      month: 'short',
      year: 'numeric'
    }).format(date);
  } catch (error) {
    console.error('Error formatting date:', dateString, error);
    return 'Неизвестно';
  }
}

function ArticleListPage() {
  const [articles, setArticles] = useState<Article[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { token, user } = useAuth();

  useEffect(() => {
    let ignore = false;

    async function loadArticles() {
      try {
        setLoading(true);
        const data = await fetchArticles(token || undefined);
        if (!ignore) {
          setArticles(Array.isArray(data) ? data : []);
          setError(null);
        }
      } catch (err) {
        if (!ignore) {
          setError(err instanceof Error ? err.message : 'Не удалось загрузить статьи');
        }
      } finally {
        if (!ignore) {
          setLoading(false);
        }
      }
    }

    loadArticles();
    return () => {
      ignore = true;
    };
  }, [token]);

  return (
    <>
      <section className="floating-panel glass" style={{ marginBottom: '2.5rem' }}>
        <h1 className="page-title">Свежие мысли сообщества</h1>
        <p className="helper-text" style={{ fontSize: '1rem', marginBottom: '1.5rem' }}>
          Подборка историй, заметок и идей, написанных авторами блога. Присоединяйтесь и делитесь своим опытом.
        </p>
        <div style={{ display: 'flex', gap: '0.75rem', flexWrap: 'wrap' }}>
          {user ? (
            <Link to="/article/new" className="primary-button">
              Написать статью
            </Link>
          ) : (
            <Link to="/register" className="primary-button">
              Присоединиться к авторам
            </Link>
          )}
          <Link to="/profile" className="secondary-button">
            Мой профиль
          </Link>
        </div>
      </section>

      {loading && <p className="helper-text">Загружаем подборку статей…</p>}

      {error && !loading && (
        <div className="empty-state">
          <p>{error}</p>
          <p>Попробуйте обновить страницу немного позже.</p>
        </div>
      )}

      {!loading && !error && articles.length === 0 && (
        <div className="empty-state">
          <p>Пока здесь тихо. Создайте первую статью и поделитесь идеями.</p>
        </div>
      )}

      {!loading && !error && articles.length > 0 && (
        <section className="card-grid">
          {articles.map((article: Article, index: number) => (
            <Link
              key={article.id}
              to={`/article/${article.id}`}
              className="card"
              style={{ animationDelay: `${index * 0.04}s` }}
            >
              <div className="chip-stack">
                <span className="chip">{formatDate(article.updated_at)}</span>
                {article.author && <span className="chip">{article.author}</span>}
                {article.viewerLiked && <span className="chip">В избранном</span>}
              </div>
              <div>
                <h2 className="card-title">{article.title}</h2>
                <p className="helper-text" style={{ marginTop: '0.35rem' }}>
                  {article.author ? `Автор: ${article.author}` : 'Автор неизвестен'} · {article.views} просмотров
                </p>
              </div>
              <p className="card-meta">
                <span className="badge">{article.likes} ❤</span>
                <span>{article.viewerLiked ? 'Нравится вам' : 'Добавьте в избранное'}</span>
              </p>
              <p className="helper-text" style={{ margin: 0, fontSize: '0.95rem', lineHeight: 1.7 }}>
                {article.content.slice(0, 220)}
                {article.content.length > 220 ? '…' : ''}
              </p>
            </Link>
          ))}
        </section>
      )}
    </>
  );
}

export default ArticleListPage;
