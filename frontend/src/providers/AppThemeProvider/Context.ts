import { createContext, useContext } from 'react'
import type {ThemeContextType} from './types';
export const ThemeContext = createContext<ThemeContextType | undefined>(
    undefined
)

export const useTheme = () => {
    const ctx = useContext(ThemeContext)
    if (!ctx) {
        throw new Error('useTheme must be used inside ThemeProvider')
    }
    return ctx
}
