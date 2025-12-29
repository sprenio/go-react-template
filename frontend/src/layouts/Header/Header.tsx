import {useState} from 'react';
import {useConfig} from '@/providers/ConfigProvider';
import {Logo} from '@/components/Logo';
import {ThemeToggle} from '@/components/ThemeToggle';
import {useTheme as useAppTheme} from '@/providers/AppThemeProvider';
import {useTranslation} from 'react-i18next';
import {Link} from 'react-router-dom';
import {useAuth, isLoggedIn} from '@/providers/AuthProvider';
import {logoutReasonCodes} from '@/api';
import {LanguageSelect} from '@/components/LanguageSelect';
import {useLang} from '@/providers/LangProvider';

export function Header() {
    const {config} = useConfig();
    const {themeName} = useAppTheme();
    const {logout} = useAuth();
    const {t} = useTranslation();

    const [menuOpen, setMenuOpen] = useState(false);
    const {lang, changeLang} = useLang();
    const toggleMenu = () => setMenuOpen((prev) => !prev);

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
                {isLoggedIn() && (
                    <>
                        <nav className="items-center gap-4 hidden md:flex" role="navigation">
                            <Link className="menu-link desktop" to="/dashboard">{t('dashboard.menu_title')}</Link>
                            <Link className="menu-link desktop" to="/settings">{t('settings.menu_title')}</Link>
                            <button className="menu-link desktop" onClick={() => logout(logoutReasonCodes.USER_LOGGED_OUT)}>
                                {t('global.logout')}
                            </button>
                        </nav>
                        <button className="block md:hidden bg-none border-none text-[var(--text)] text-[length:var(--text-xxl)] cursor-pointer" onClick={toggleMenu} aria-label="Menu">
                            ☰
                        </button>
                        {menuOpen && (
                            <nav className={`flex md:hidden flex-col bg-[var(--surface)] border-t border-[var(--border)] shadow-md absolute top-full right-0 w-full overflow-hidden z-[999]
    transition-[max-height,opacity] duration-300 ease-in-out max-h-[200px]`}>
                                <Link className="menu-link mobile" to="/dashboard" onClick={() => setMenuOpen(false)}>
                                    {t('dashboard.menu_title')}
                                </Link>
                                <Link className="menu-link mobile" to="/settings" onClick={() => setMenuOpen(false)}>
                                    {t('settings.menu_title')}
                                </Link>
                                <button className="menu-link mobile" onClick={() => logout(logoutReasonCodes.USER_LOGGED_OUT)}>
                                    {t('global.logout')}
                                </button>
                            </nav>)
                        }
                    </>
                )}
                <LanguageSelect value={lang} onChange={handleLangChange} flagsOnly={true}/>
                <ThemeToggle/>
            </div>
        </header>
    );
};

