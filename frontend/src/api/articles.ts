import type { Article, ArticlePayload, LikeResponse } from '../types';

const BASE_URL = '/api/articles';

async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const message = await extractError(response);
    throw new Error(message);
  }
  if (response.status === 204) {
    return undefined as T;
  }
  return response.json() as Promise<T>;
}

async function extractError(response: Response): Promise<string> {
  try {
    const data = await response.json();
    if (data?.error) {
      return data.error as string;
    }
  } catch (error) {
    // Ignore JSON parse issues and fall back to status text
  }
  return response.statusText || 'Unknown error';
}

function withAuth(options: RequestInit = {}, token?: string): RequestInit {
  const headers = new Headers(options.headers);
  if (options.body && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json');
  }
  if (token) {
    headers.set('Authorization', `Bearer ${token}`);
  }

  return {
    ...options,
    headers
  };
}

export async function fetchArticles(token?: string): Promise<Article[]> {
  const response = await fetch(BASE_URL, withAuth(undefined, token));
  return handleResponse<Article[]>(response);
}

export async function fetchArticle(id: number, token?: string): Promise<Article> {
  const response = await fetch(`${BASE_URL}/${id}`, withAuth(undefined, token));
  return handleResponse<Article>(response);
}

export async function fetchArticlePreview(id: number, token: string): Promise<Article> {
  const response = await fetch(`${BASE_URL}/${id}?preview=true`, withAuth(undefined, token));
  return handleResponse<Article>(response);
}

export async function createArticle(payload: ArticlePayload, token: string): Promise<Article> {
  const response = await fetch(
    BASE_URL,
    withAuth(
      {
        method: 'POST',
        body: JSON.stringify(payload)
      },
      token
    )
  );
  return handleResponse<Article>(response);
}

export async function updateArticle(id: number, payload: ArticlePayload, token: string): Promise<Article> {
  const response = await fetch(
    `${BASE_URL}/${id}`,
    withAuth(
      {
        method: 'PUT',
        body: JSON.stringify(payload)
      },
      token
    )
  );
  return handleResponse<Article>(response);
}

export async function deleteArticle(id: number, token: string): Promise<void> {
  const response = await fetch(`${BASE_URL}/${id}`, withAuth({ method: 'DELETE' }, token));
  await handleResponse<void>(response);
}

export async function likeArticle(id: number, token: string, like = true): Promise<LikeResponse> {
  const response = await fetch(
    `${BASE_URL}/${id}/like`,
    withAuth(
      {
        method: 'POST',
        body: JSON.stringify({ like })
      },
      token
    )
  );
  return handleResponse<LikeResponse>(response);
}
