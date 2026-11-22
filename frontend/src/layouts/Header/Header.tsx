import { useConfig } from '@/providers/ConfigProvider';
import { Logo } from '@/components/Logo';
import { ThemeToggle } from '@/components/ThemeToggle';
import { useTheme as useAppTheme } from '@/providers/AppThemeProvider';
import { LanguageSelect } from '@/components/LanguageSelect';
import { useLang } from '@/providers/LangProvider';
import { TopMenu } from '@/layouts/TopMenu';

export function Header() {
  const { config } = useConfig();
  const { themeName } = useAppTheme();
  const { lang, changeLang } = useLang();
  const handleLangChange = (langCode: string): void => {
    changeLang(langCode);
  };
  return (
    <header className="relative z-[1000] flex items-center justify-between border-b border-[var(--border)] bg-[var(--surface)] px-6 py-3 text-[var(--text)] shadow-sm">
      <div className="flex items-center gap-3">
        <h1 className="m-0 flex cursor-pointer items-center gap-2 text-[length:var(--text-lg)] font-medium transition hover:-translate-y-[1px] hover:text-[var(--primary)]">
          <span className="h-9 w-15 transition-transform hover:scale-105 hover:rotate-[-5deg]">
            <Logo theme={themeName} />
          </span>
          {config.AppName}
        </h1>
      </div>

      <div className="top-menu flex items-center gap-4">
        <TopMenu
          menuItems={[
            'dashboard',
            'settings',
          ]}
        />
        <LanguageSelect value={lang} onChange={handleLangChange} flagsOnly={true} />
        <ThemeToggle />
      </div>
    </header>
  );
}
