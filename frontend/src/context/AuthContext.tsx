import { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react';
import type { LoginPayload, ProfileUpdatePayload, RegisterPayload, User } from '../types';
import * as authApi from '../api/auth';

interface AuthContextValue {
  user: User | null;
  token: string | null;
  loading: boolean;
  error: string | null;
  login: (payload: LoginPayload) => Promise<void>;
  register: (payload: RegisterPayload) => Promise<void>;
  updateProfile: (payload: ProfileUpdatePayload) => Promise<void>;
  logout: () => void;
  refreshProfile: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);
const STORAGE_KEY = 'blog::token';

function normalizeUser(authUser: User): User {
  const contacts = Array.isArray((authUser as unknown as { contacts?: unknown }).contacts)
    ? ((authUser as unknown as { contacts: unknown[] }).contacts.filter((value): value is string => typeof value === 'string'))
    : undefined;

  return {
    ...authUser,
    contacts
  };
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(() => localStorage.getItem(STORAGE_KEY));
  const [loading, setLoading] = useState<boolean>(!!token);
  const [error, setError] = useState<string | null>(null);

  const applyAuth = useCallback((authToken: string, authUser: User) => {
    localStorage.setItem(STORAGE_KEY, authToken);
    setToken(authToken);
    setUser(normalizeUser(authUser));
    setError(null);
  }, []);

  const clearAuth = useCallback(() => {
    localStorage.removeItem(STORAGE_KEY);
    setToken(null);
    setUser(null);
    setError(null);
  }, []);

  const refreshProfile = useCallback(async () => {
    if (!token) {
        setUser(null);
        setLoading(false);
      return;
    }

    try {
      setLoading(true);
      console.log('[AuthContext] Fetching profile with token:', token.substring(0, 20) + '...');
      const profile = await authApi.fetchProfile(token);
      console.log('[AuthContext] Profile fetched successfully:', profile);
  setUser(normalizeUser(profile));
      setError(null);
    } catch (err) {
      console.error('[AuthContext] Failed to fetch profile:', err);
      clearAuth();
      setError(err instanceof Error ? err.message : 'Не удалось получить профиль');
    } finally {
      setLoading(false);
    }
  }, [token, clearAuth]);

  useEffect(() => {
    if (token) {
      refreshProfile();
    }
  }, [token, refreshProfile]);

  const login = useCallback(
    async (payload: LoginPayload) => {
      try {
        setLoading(true);
        console.log('[AuthContext] Attempting login...');
        const { token: authToken, user: authUser } = await authApi.login(payload);
        console.log('[AuthContext] Login successful, token:', authToken.substring(0, 20) + '...', 'user:', authUser);
  applyAuth(authToken, authUser);
        setError(null);
      } catch (err) {
        console.error('[AuthContext] Login failed:', err);
        const message = err instanceof Error ? err.message : 'Не удалось выполнить вход';
        setError(message);
        throw err;
      } finally {
        setLoading(false);
      }
    },
    [applyAuth]
  );

  const register = useCallback(
    async (payload: RegisterPayload) => {
      try {
        setLoading(true);
        const { token: authToken, user: authUser } = await authApi.register(payload);
  applyAuth(authToken, authUser);
        setError(null);
      } catch (err) {
        const message = err instanceof Error ? err.message : 'Не удалось завершить регистрацию';
        setError(message);
        throw err;
      } finally {
        setLoading(false);
      }
    },
    [applyAuth]
  );

  const logout = useCallback(() => {
    clearAuth();
  }, [clearAuth]);

  const updateProfile = useCallback(
    async (payload: ProfileUpdatePayload) => {
      if (!token) {
        throw new Error('Требуется авторизация');
      }

      try {
        setLoading(true);
        const updated = await authApi.updateProfile(token, payload);
  setUser(normalizeUser(updated));
        setError(null);
      } catch (err) {
        const message = err instanceof Error ? err.message : 'Не удалось обновить профиль';
        setError(message);
        throw err;
      } finally {
        setLoading(false);
      }
    },
    [token]
  );

  const value = useMemo<AuthContextValue>(
    () => ({ user, token, loading, error, login, register, updateProfile, logout, refreshProfile }),
    [user, token, loading, error, login, register, updateProfile, logout, refreshProfile]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextValue {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
