import { createContext, useContext } from 'react'
import type {LangContextType} from './types';
export const LangContext = createContext<LangContextType | undefined>(
    undefined
)

export const useLang = () => {
    const ctx = useContext(LangContext)
    if (!ctx) {
        throw new Error('useLang must be used inside LangProvider')
    }
    return ctx
}
