import {useConfig} from '@/providers/ConfigProvider';
import {Logo} from '@/components/Logo';
import {ThemeToggle} from '@/components/ThemeToggle';
import {useTheme as useAppTheme} from '@/providers/AppThemeProvider';
import {LanguageSelect} from '@/components/LanguageSelect';
import {useLang} from '@/providers/LangProvider';
import { TopMenu } from '@/layouts/TopMenu';

export function Header() {
    const {config} = useConfig();
    const {themeName} = useAppTheme();
    const {lang, changeLang} = useLang();
    const handleLangChange = (langCode: string): void => {
        changeLang(langCode);
    };
    return (
        <header className="relative z-[1000] flex items-center justify-between px-6 py-3 bg-[var(--surface)] text-[var(--text)] border-b border-[var(--border)] shadow-sm">
            <div className="flex items-center gap-3">
                <h1 className="flex items-center gap-2 m-0 text-[length:var(--text-lg)] font-medium cursor-pointer transition hover:text-[var(--primary)] hover:-translate-y-[1px] ">
                    <span className="w-15 h-9 transition-transform hover:rotate-[-5deg] hover:scale-105 " >
                        <Logo theme={themeName}/>
                    </span>
                    {config.AppName}
                </h1>
            </div>

            <div className="flex items-center gap-4 top-menu">
                <TopMenu/>
                <LanguageSelect value={lang} onChange={handleLangChange} flagsOnly={true}/>
                <ThemeToggle/>
            </div>
        </header>
    );
}
