export type ThemeNameType = 'light' | 'dark';

export type ThemeContextType = {
    themeName: ThemeNameType
    toggleTheme: () => void
}
