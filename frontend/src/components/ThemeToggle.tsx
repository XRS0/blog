import { useTheme } from '../context/ThemeContext';

export function ThemeToggle() {
  const { themeName, toggleTheme, theme } = useTheme();
  const isDark = themeName === 'dark';

  return (
    <button
      type="button"
      className="theme-toggle"
      onClick={toggleTheme}
      aria-label={`Переключить тему. Сейчас ${theme.displayName.toLowerCase()} тема.`}
    >
      <span aria-hidden="true" style={{ fontSize: '1.05rem' }}>
        {isDark ? '🌙' : '☀️'}
      </span>
      <span>{isDark ? 'Тёмная' : 'Светлая'}</span>
    </button>
  );
}
