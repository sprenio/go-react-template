import { createContext, useContext } from 'react'
import type {ConfigContextType} from './types';
export const ConfigContext = createContext<ConfigContextType | undefined>(
    undefined
)

export const useConfig = () => {
    const ctx = useContext(ConfigContext)
    if (!ctx) {
        throw new Error('useConfig must be used inside ConfigProvider')
    }
    return ctx
}
