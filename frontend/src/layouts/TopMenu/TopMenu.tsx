import {useState} from 'react';
import { Link } from 'react-router-dom';
import { logoutReasonCodes } from '@/api';
import { useAuth } from '@/providers/AuthProvider';
import { useTranslation } from 'react-i18next';

export function TopMenu() {
  const {logout, appUser} = useAuth();
  const {t} = useTranslation();
  const [menuOpen, setMenuOpen] = useState(false);
  const toggleMenu = () => setMenuOpen((prev) => !prev);
  if(!appUser) {
    return null
  }
  const menuItems = ['dashboard', 'settings'];
  return (
    <>
      <nav className="items-center gap-4 hidden md:flex" role="navigation">
        {menuItems.map((value, index) => (
          <Link key={index} className="menu-link desktop" to={"/"+value}>
            {t(value+'.menu_title')}
          </Link>
        ))}
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
          {menuItems.map((value, index) => (
            <Link key={index} className="menu-link mobile" to={"/"+value} onClick={() => setMenuOpen(false)}>
              {t(value+'.menu_title')}
            </Link>
          ))}
          <button className="menu-link mobile" onClick={() => logout(logoutReasonCodes.USER_LOGGED_OUT)}>
            {t('global.logout')}
          </button>
        </nav>)
      }
    </>
  );
}
