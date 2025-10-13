export type ThemeName = 'light' | 'dark';

export type ThemeTokens = Record<string, string>;

export interface ThemeDefinition {
  name: ThemeName;
  displayName: string;
  tokens: ThemeTokens;
}

export const ACCENT_COLOR = '#3b82f6';

export const lightTheme: ThemeDefinition = {
  name: 'light',
  displayName: 'Светлая',
  tokens: {
    '--color-accent': ACCENT_COLOR,
    '--color-accent-soft': 'rgba(59, 130, 246, 0.12)',
    '--color-background': '#f0f9ff',
    '--color-background-alt': '#ffffff',
    '--color-surface': '#ffffff',
    '--color-surface-elevated': '#f0f9ff',
    '--color-border': 'rgba(30, 64, 175, 0.16)',
    '--color-text-primary': '#0f172a',
    '--color-text-secondary': 'rgba(15, 23, 42, 0.7)',
    '--color-text-muted': 'rgba(15, 23, 42, 0.55)',
    '--color-shadow': 'rgba(59, 130, 246, 0.18)',
    '--color-shadow-strong': 'rgba(59, 130, 246, 0.32)',
    '--color-glow': 'rgba(59, 130, 246, 0.18)',
    '--gradient-background': 'radial-gradient(circle at top left, rgba(59, 130, 246, 0.16), transparent 55%), radial-gradient(circle at bottom right, rgba(59, 130, 246, 0.12), transparent 50%)',
  }
};

export const darkTheme: ThemeDefinition = {
  name: 'dark',
  displayName: 'Тёмная',
  tokens: {
    '--color-accent': ACCENT_COLOR,
    '--color-accent-soft': 'rgba(68, 136, 247, 0.25)',
    '--color-background': '#18243aff',
    '--color-background-alt': '#111827',
    '--color-surface': '#212e43ff',
    '--color-surface-elevated': '#334155',
    '--color-border': 'rgba(96, 165, 250, 0.28)',
    '--color-text-primary': '#f1f5f9',
    '--color-text-secondary': 'rgba(241, 245, 249, 0.78)',
    '--color-text-muted': 'rgba(241, 245, 249, 0.62)',
    '--color-shadow': 'rgba(59, 131, 246, 0.09)',
    '--color-shadow-strong': 'rgba(59, 130, 246, 0.38)',
    '--color-glow': 'rgba(59, 130, 246, 0.32)',
    '--gradient-background': 'radial-gradient(circle at top left, rgba(59, 130, 246, 0.22), transparent 55%), radial-gradient(circle at bottom right, rgba(59, 130, 246, 0.18), transparent 50%)',
  }
};

export const themes: Record<ThemeName, ThemeDefinition> = {
  light: lightTheme,
  dark: darkTheme,
};

export const orderedThemes: ThemeDefinition[] = [lightTheme, darkTheme];
