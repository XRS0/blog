# Frontend Integration Guide

## Обновления API для поддержки visibility

### Изменения в типах

```typescript
// types.ts
export type ArticleVisibility = 'public' | 'private' | 'link';

export interface Article {
  id: number;
  title: string;
  content: string;
  visibility: ArticleVisibility;
  access_token?: string;  // Только для visibility='link'
  access_url?: string;    // Полная ссылка для шаринга
  author: string;
  views: number;
  likes: number;
  viewerLiked: boolean;
  created_at: string;
  updated_at: string;
}
```

### Обновления в api/articles.ts

```typescript
export interface CreateArticleRequest {
  title: string;
  content: string;
  visibility: ArticleVisibility;
}

export interface UpdateArticleRequest {
  title: string;
  content: string;
  visibility: ArticleVisibility;
}

export const createArticle = async (data: CreateArticleRequest): Promise<Article> => {
  const response = await api.post('/api/articles', data);
  return response.data;
};

export const updateArticle = async (id: number, data: UpdateArticleRequest): Promise<Article> => {
  const response = await api.put(`/api/articles/${id}`, data);
  return response.data;
};

// Новый метод для получения статьи с токеном доступа
export const getArticleWithToken = async (id: number, accessToken: string): Promise<Article> => {
  const response = await api.get(`/api/articles/${id}?access_token=${accessToken}`);
  return response.data;
};
```

### Компонент выбора видимости

```typescript
// components/VisibilitySelector.tsx
import { useState } from 'react';

interface VisibilitySelectorProps {
  value: ArticleVisibility;
  onChange: (value: ArticleVisibility) => void;
}

export function VisibilitySelector({ value, onChange }: VisibilitySelectorProps) {
  return (
    <div className="visibility-selector">
      <label>
        <input
          type="radio"
          value="public"
          checked={value === 'public'}
          onChange={(e) => onChange(e.target.value as ArticleVisibility)}
        />
        <span>📢 Public</span>
        <small>Видно всем пользователям</small>
      </label>

      <label>
        <input
          type="radio"
          value="private"
          checked={value === 'private'}
          onChange={(e) => onChange(e.target.value as ArticleVisibility)}
        />
        <span>🔒 Private</span>
        <small>Только вы можете видеть</small>
      </label>

      <label>
        <input
          type="radio"
          value="link"
          checked={value === 'link'}
          onChange={(e) => onChange(e.target.value as ArticleVisibility)}
        />
        <span>🔗 Link</span>
        <small>Доступ по специальной ссылке</small>
      </label>
    </div>
  );
}
```

### Обновление ArticleEditorPage

```typescript
// pages/ArticleEditorPage.tsx
import { VisibilitySelector } from '../components/VisibilitySelector';

export function ArticleEditorPage() {
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [visibility, setVisibility] = useState<ArticleVisibility>('public');
  const [accessUrl, setAccessUrl] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    try {
      const article = await createArticle({ title, content, visibility });
      
      // Если создана статья с visibility='link', показать ссылку
      if (article.visibility === 'link' && article.access_url) {
        setAccessUrl(article.access_url);
      }
      
      navigate(`/articles/${article.id}`);
    } catch (error) {
      console.error('Failed to create article:', error);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="text"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        placeholder="Article Title"
      />
      
      <textarea
        value={content}
        onChange={(e) => setContent(e.target.value)}
        placeholder="Write your article..."
      />

      <VisibilitySelector value={visibility} onChange={setVisibility} />

      {accessUrl && (
        <div className="access-url-notice">
          <p>✅ Article created! Share this link:</p>
          <input type="text" value={window.location.origin + accessUrl} readOnly />
          <button onClick={() => navigator.clipboard.writeText(window.location.origin + accessUrl)}>
            Copy Link
          </button>
        </div>
      )}

      <button type="submit">Publish</button>
    </form>
  );
}
```

### Просмотр статьи с access token

```typescript
// pages/ArticleDetailPage.tsx
import { useParams, useSearchParams } from 'react-router-dom';

export function ArticleDetailPage() {
  const { id } = useParams();
  const [searchParams] = useSearchParams();
  const accessToken = searchParams.get('access_token');
  
  const [article, setArticle] = useState<Article | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchArticle = async () => {
      try {
        let data;
        if (accessToken) {
          data = await getArticleWithToken(Number(id), accessToken);
        } else {
          data = await getArticle(Number(id));
        }
        setArticle(data);
      } catch (err: any) {
        if (err.response?.status === 403) {
          setError('Access denied. This article is private or requires an access token.');
        } else {
          setError('Failed to load article');
        }
      }
    };

    fetchArticle();
  }, [id, accessToken]);

  if (error) {
    return <div className="error">{error}</div>;
  }

  if (!article) {
    return <div>Loading...</div>;
  }

  return (
    <article>
      <header>
        <h1>{article.title}</h1>
        <div className="meta">
          <span>By {article.author}</span>
          <span>{article.views} views</span>
          <span>{article.likes} likes</span>
          
          {article.visibility === 'private' && (
            <span className="badge private">🔒 Private</span>
          )}
          {article.visibility === 'link' && (
            <span className="badge link">🔗 Shared via link</span>
          )}
        </div>
      </header>
      
      <MarkdownRenderer content={article.content} />
      
      {article.visibility === 'link' && article.access_url && (
        <div className="share-section">
          <h3>Share this article</h3>
          <input 
            type="text" 
            value={window.location.href} 
            readOnly 
          />
          <button onClick={() => navigator.clipboard.writeText(window.location.href)}>
            Copy Link
          </button>
        </div>
      )}
    </article>
  );
}
```

### Фильтрация статей в профиле

```typescript
// pages/ProfilePage.tsx
export function ProfilePage() {
  const { user } = useAuth();
  const [articles, setArticles] = useState<Article[]>([]);
  const [filter, setFilter] = useState<'all' | ArticleVisibility>('all');

  const filteredArticles = useMemo(() => {
    if (filter === 'all') return articles;
    return articles.filter(a => a.visibility === filter);
  }, [articles, filter]);

  return (
    <div>
      <h1>My Articles</h1>
      
      <div className="filters">
        <button onClick={() => setFilter('all')}>All</button>
        <button onClick={() => setFilter('public')}>📢 Public</button>
        <button onClick={() => setFilter('private')}>🔒 Private</button>
        <button onClick={() => setFilter('link')}>🔗 Link</button>
      </div>

      <div className="articles-grid">
        {filteredArticles.map(article => (
          <ArticleCard key={article.id} article={article} />
        ))}
      </div>
    </div>
  );
}
```

### Индикаторы видимости в карточке статьи

```typescript
// components/ArticleCard.tsx
export function ArticleCard({ article }: { article: Article }) {
  return (
    <div className="article-card">
      <h3>{article.title}</h3>
      
      <div className="visibility-badge">
        {article.visibility === 'public' && '📢 Public'}
        {article.visibility === 'private' && '🔒 Private'}
        {article.visibility === 'link' && '🔗 Link only'}
      </div>
      
      <p>{article.content.substring(0, 150)}...</p>
      
      {article.visibility === 'link' && article.access_url && (
        <button onClick={() => {
          navigator.clipboard.writeText(
            window.location.origin + article.access_url
          );
          alert('Link copied!');
        }}>
          📋 Copy Link
        </button>
      )}
      
      <Link to={`/articles/${article.id}`}>Read more</Link>
    </div>
  );
}
```

## Стили (пример)

```css
.visibility-selector {
  display: flex;
  gap: 1rem;
  margin: 1rem 0;
}

.visibility-selector label {
  display: flex;
  flex-direction: column;
  padding: 1rem;
  border: 2px solid #ddd;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.visibility-selector label:has(input:checked) {
  border-color: #007bff;
  background-color: #f0f8ff;
}

.visibility-badge {
  display: inline-block;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.875rem;
  margin: 0.5rem 0;
}

.badge.private {
  background-color: #ff4444;
  color: white;
}

.badge.link {
  background-color: #ffa500;
  color: white;
}

.access-url-notice {
  padding: 1rem;
  background-color: #d4edda;
  border: 1px solid #c3e6cb;
  border-radius: 4px;
  margin: 1rem 0;
}

.share-section {
  margin-top: 2rem;
  padding: 1rem;
  background-color: #f8f9fa;
  border-radius: 8px;
}

.share-section input {
  width: 100%;
  padding: 0.5rem;
  margin: 0.5rem 0;
  border: 1px solid #ddd;
  border-radius: 4px;
}
```

## Тестирование

1. **Public article**: Создайте публичную статью, проверьте что она видна в списке всех статей
2. **Private article**: Создайте приватную статью, проверьте что она видна только вам
3. **Link article**: Создайте статью с visibility='link', скопируйте ссылку, откройте в инкогнито - должна открыться
4. **Access denied**: Попробуйте открыть link-статью без токена - должен быть 403
5. **Update visibility**: Измените visibility существующей статьи

## Миграция существующих статей

Все существующие статьи автоматически получат `visibility='public'` после применения SQL миграции.
