import type { AuthResponse, LoginPayload, ProfileUpdatePayload, RegisterPayload, User } from '../types';

const BASE_URL = '/api/auth';

async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const message = await extractError(response);
    throw new Error(message);
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

export async function register(payload: RegisterPayload): Promise<AuthResponse> {
  const response = await fetch(`${BASE_URL}/register`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(payload)
  });
  return handleResponse<AuthResponse>(response);
}

export async function login(payload: LoginPayload): Promise<AuthResponse> {
  const response = await fetch(`${BASE_URL}/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(payload)
  });
  return handleResponse<AuthResponse>(response);
}

export async function fetchProfile(token: string): Promise<User> {
  const response = await fetch(`${BASE_URL}/me`, {
    headers: {
      Authorization: `Bearer ${token}`
    }
  });
  const data = await handleResponse<{ user: User }>(response);
  return data.user;
}

export async function updateProfile(token: string, payload: ProfileUpdatePayload): Promise<User> {
  const response = await fetch(`${BASE_URL}/me`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`
    },
    body: JSON.stringify(payload)
  });
  const data = await handleResponse<{ user: User }>(response);
  return data.user;
}
