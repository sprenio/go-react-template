import React, { useState } from 'react';
import i18n from '@/i18n';
import { LangContext } from './Context';


export function LangProvider({ children }: { children: React.ReactNode }) {
    const [lang, setLang] = useState<string>(localStorage.getItem('userLanguage') || 'pl');
    const changeLang = (newLang:string) => {
        setLang(newLang);
        localStorage.setItem('userLanguage', newLang);
        i18n.changeLanguage(newLang);
    };
    return <LangContext.Provider value={{ lang, changeLang }}>{children}</LangContext.Provider>;
}
