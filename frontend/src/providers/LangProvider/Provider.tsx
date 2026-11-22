import React, { useState } from 'react';
import i18n from '@/i18n';
import { LangContext } from './Context';
import { registerLocale } from 'react-datepicker';
import { pl, enGB, de, uk, type Locale } from 'date-fns/locale';

const localeCodeMap: Record<string, string> = {
  pl: 'pl',
  en: 'en-GB',
  de: 'de',
  uk: 'uk',
};
const localeObjectMap: Record<string, Locale> = {
  pl,
  'en-GB': enGB,
  de,
  uk,
};

Object.entries(localeObjectMap).forEach(([code, locale]) => {
  registerLocale(code, locale);
});

function resolveLocale(lang: string): string {
  return localeCodeMap[lang] || 'en-GB';
}

export function LangProvider({ children }: { children: React.ReactNode }) {
  const [lang, setLang] = useState<string>(localStorage.getItem('userLanguage') || 'pl');
  const [locale, setLocale] = useState<string>(resolveLocale(lang));
  const changeLang = (newLang: string) => {
    setLang(newLang);
    setLocale(resolveLocale(newLang));
    localStorage.setItem('userLanguage', newLang);
    i18n.changeLanguage(newLang);
  };
  return (
    <LangContext.Provider value={{ lang, locale, changeLang }}>{children}</LangContext.Provider>
  );
}
