export interface User {
  id: number;
  email: string;
  username: string;
  contacts?: string[];
  createdAt?: string;
  updatedAt?: string;
}

export interface Article {
  id: number;
  userId: number;
  title: string;
  content: string;
  views: number;
  likes: number;
  viewerLiked: boolean;
  author: string;
  created_at: string;
  updated_at: string;
}

export type ArticlePayload = Pick<Article, 'title' | 'content'>;

export interface LoginPayload {
  email: string;
  password: string;
}

export interface RegisterPayload extends LoginPayload {
  username: string;
  contacts: string[];
}

export interface ProfileUpdatePayload {
  username: string;
  contacts: string[];
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface LikeResponse {
  likes: number;
  success: boolean;
}
