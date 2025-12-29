import React, { useState, useEffect } from 'react';
import { ThemeContext } from './Context';
import type {ThemeNameType} from './types';

const THEME_DARK = 'dark';
const THEME_LIGHT = 'light';

export function AppThemeProvider({ children }: { children: React.ReactNode }) {
    const getInitialTheme = ():ThemeNameType => {
        const saved:string = localStorage.getItem('theme') || '';
        if (saved) {
            return saved === THEME_DARK ? THEME_DARK : THEME_LIGHT;
        }
        const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
        return prefersDark ? THEME_DARK : THEME_LIGHT;
    };
    const [themeName, setThemeName] = useState(getInitialTheme);

    useEffect(() => {
        const root = document.documentElement; // <html />
        root.setAttribute('data-theme', themeName);
        localStorage.setItem('theme', themeName);
    }, [themeName]);

    const toggleTheme = () => {
        setThemeName((prev) => (prev === THEME_LIGHT ? THEME_DARK : THEME_LIGHT));
    };

    return <ThemeContext.Provider value={{ themeName, toggleTheme }}>{children}</ThemeContext.Provider>;
}
