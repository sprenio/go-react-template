import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';

// PL
import plBase from '@/locales/pl/base.json';
import plLogin from '@/locales/pl/login.json';
import plDashboard from '@/locales/pl/dashboard.json';
import plSettings from '@/locales/pl/settings.json';
import plRegister from '@/locales/pl/register.json';
import plReset from '@/locales/pl/reset.json';
import plConfirm from '@/locales/pl/confirm.json';
import plPasswordChange from '@/locales/pl/password_change.json';

// EN
import enBase from '@/locales/en/base.json';
import enLogin from '@/locales/en/login.json';
import enDashboard from '@/locales/en/dashboard.json';
import enSettings from '@/locales/en/settings.json';
import enRegister from '@/locales/en/register.json';
import enReset from '@/locales/en/reset.json';
import enConfirm from '@/locales/en/confirm.json';
import enPasswordChange from '@/locales/en/password_change.json';

// DE
import deBase from '@/locales/de/base.json';
import deLogin from '@/locales/de/login.json';
import deDashboard from '@/locales/de/dashboard.json';
import deSettings from '@/locales/de/settings.json';
import deRegister from '@/locales/de/register.json';
import deReset from '@/locales/de/reset.json';
import deConfirm from '@/locales/de/confirm.json';
import dePasswordChange from '@/locales/de/password_change.json';

// UK
import ukBase from '@/locales/uk/base.json';
import ukLogin from '@/locales/uk/login.json';
import ukDashboard from '@/locales/uk/dashboard.json';
import ukSettings from '@/locales/uk/settings.json';
import ukRegister from '@/locales/uk/register.json';
import ukReset from '@/locales/uk/reset.json';
import ukConfirm from '@/locales/uk/confirm.json';
import ukPasswordChange from '@/locales/uk/password_change.json';

function pluralizePl(count: number, one: string, few: string, many: string): string {
  const abs = Math.abs(count);
  if (abs === 1) return one;
  if (abs % 10 >= 2 && abs % 10 <= 4 && (abs % 100 < 12 || abs % 100 > 14)) {
    return few;
  }
  return many;
}

const userLang = localStorage.getItem('userLanguage') || 'pl';

i18n.use(initReactI18next).init({
  resources: {
    pl: {
      translation: {
        ...plBase,
        login: plLogin,
        dashboard: plDashboard,
        settings: plSettings,
        register: plRegister,
        reset: plReset,
        confirm: plConfirm,
        password_change: plPasswordChange,
      },
    },
    en: {
      translation: {
        ...enBase,
        login: enLogin,
        dashboard: enDashboard,
        settings: enSettings,
        register: enRegister,
        reset: enReset,
        confirm: enConfirm,
        password_change: enPasswordChange,
      },
    },
    de: {
      translation: {
        ...deBase,
        login: deLogin,
        dashboard: deDashboard,
        settings: deSettings,
        register: deRegister,
        reset: deReset,
        confirm: deConfirm,
        password_change: dePasswordChange,
      },
    },
    uk: {
      translation: {
        ...ukBase,
        login: ukLogin,
        dashboard: ukDashboard,
        settings: ukSettings,
        register: ukRegister,
        reset: ukReset,
        confirm: ukConfirm,
        password_change: ukPasswordChange,
      },
    },
  },
  lng: userLang,
  fallbackLng: 'en',
  interpolation: {
    escapeValue: false,
    format: (value, format) => {
      if (format && format.startsWith('pluralPl')) {
        const match = format.match(/\(([^)]+)\)/);
        if (match) {
          const [one, few, many] = match[1].split(',').map((s) => s.trim());
          return pluralizePl(value, one, few, many);
        }
      }
      return value;
    },
  },
});

export default i18n;
