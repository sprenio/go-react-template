import { useState } from 'react';
import { Link } from 'react-router-dom';
import { logoutReasonCodes } from '@/api';
import { useAuth } from '@/providers/AuthProvider';
import { useTranslation } from 'react-i18next';
import type { TopMenuProps } from './types';

export function TopMenu({ menuItems }: TopMenuProps) {
  const { logout, appUser } = useAuth();
  const { t } = useTranslation();
  const [menuOpen, setMenuOpen] = useState(false);
  const toggleMenu = () => setMenuOpen((prev) => !prev);
  if (!appUser) {
    return null;
  }
  //const menuItems = ['sources', 'dashboard', 'settings'];
  return (
    <>
      <nav className="hidden items-center gap-4 xl:flex" role="navigation">
        {menuItems.map((value, index) => (
          <Link key={index} className="menu-link desktop" to={'/' + value}>
            {t(value + '.menu_title')}
          </Link>
        ))}
        <button
          className="menu-link desktop"
          onClick={() => logout(logoutReasonCodes.USER_LOGGED_OUT)}
        >
          {t('global.logout')}
        </button>
      </nav>
      <button
        className="block cursor-pointer border-none bg-none text-[length:var(--text-xxl)] text-[var(--text)] xl:hidden"
        onClick={toggleMenu}
        aria-label="Menu"
      >
        ☰
      </button>
      {menuOpen && (
        <nav
          className={`absolute top-full right-0 z-[999] flex max-h-[calc(100vh-100%)] w-full flex-col overflow-y-auto border-t border-[var(--border)] bg-[var(--surface)] shadow-md xl:hidden`}
        >
          {menuItems.map((value, index) => (
            <Link
              key={index}
              className="menu-link mobile"
              to={'/' + value}
              onClick={() => setMenuOpen(false)}
            >
              {t(value + '.menu_title')}
            </Link>
          ))}
          <button
            className="menu-link mobile"
            onClick={() => logout(logoutReasonCodes.USER_LOGGED_OUT)}
          >
            {t('global.logout')}
          </button>
        </nav>
      )}
    </>
  );
}
