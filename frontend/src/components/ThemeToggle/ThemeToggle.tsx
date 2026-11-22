import { Sun, Moon } from 'lucide-react';
import { useTheme as useAppTheme } from '@/providers/AppThemeProvider';

const ThemeToggle = () => {
  const { themeName, toggleTheme } = useAppTheme();

  const rootStyles = getComputedStyle(document.documentElement);
  const sunColor = rootStyles.getPropertyValue('--icon-sun').trim();
  const moonColor = rootStyles.getPropertyValue('--icon-mon').trim();

  return (
    <button
      onClick={toggleTheme}
      aria-label="Toggle theme"
      className="flex cursor-pointer items-center justify-center rounded-[var(--radius-round)] border-none bg-transparent transition-all duration-[var(--transition-fast)] hover:bg-[var(--muted)]"
    >
      {themeName === 'light' ? (
        <Moon color={moonColor} size={22} />
      ) : (
        <Sun color={sunColor} size={24} />
      )}
    </button>
  );
};

export default ThemeToggle;
