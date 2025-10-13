import { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react';
import { ThemeDefinition, ThemeName, themes, orderedThemes } from '../theme/tokens';

interface ThemeContextValue {
  theme: ThemeDefinition;
  themeName: ThemeName;
  availableThemes: ThemeDefinition[];
  setTheme: (name: ThemeName) => void;
  toggleTheme: () => void;
}

const STORAGE_KEY = 'blog::theme';

const ThemeContext = createContext<ThemeContextValue | undefined>(undefined);

function applyThemeTokens(theme: ThemeDefinition) {
  const root = document.documentElement;
  root.setAttribute('data-theme', theme.name);
  Object.entries(theme.tokens).forEach(([key, value]) => {
    root.style.setProperty(key, value);
  });
}

function readStoredTheme(): ThemeName | null {
  if (typeof window === 'undefined') {
    return null;
  }
  const stored = window.localStorage.getItem(STORAGE_KEY);
  if (stored === 'light' || stored === 'dark') {
    return stored;
  }
  return null;
}

function detectPreferredTheme(): ThemeName {
  if (typeof window !== 'undefined' && window.matchMedia?.('(prefers-color-scheme: dark)').matches) {
    return 'dark';
  }
  return 'light';
}

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const [themeName, setThemeName] = useState<ThemeName>(() => {
    return readStoredTheme() ?? detectPreferredTheme();
  });

  useEffect(() => {
    const theme = themes[themeName];
    applyThemeTokens(theme);
    window.localStorage.setItem(STORAGE_KEY, theme.name);
  }, [themeName]);

  const setTheme = useCallback((name: ThemeName) => {
    setThemeName(name);
  }, []);

  const toggleTheme = useCallback(() => {
    setThemeName((prev) => (prev === 'light' ? 'dark' : 'light'));
  }, []);

  const value = useMemo<ThemeContextValue>(() => {
    const theme = themes[themeName];
    return {
      theme,
      themeName,
      availableThemes: orderedThemes,
      setTheme,
      toggleTheme,
    };
  }, [setTheme, themeName, toggleTheme]);

  return <ThemeContext.Provider value={value}>{children}</ThemeContext.Provider>;
}

export function useTheme(): ThemeContextValue {
  const ctx = useContext(ThemeContext);
  if (!ctx) {
    throw new Error('useTheme must be used inside ThemeProvider');
  }
  return ctx;
}
