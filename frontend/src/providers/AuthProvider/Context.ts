import { createContext, useContext } from 'react'
import type { AuthContextType } from './types'

export const AuthContext = createContext<AuthContextType | undefined>(
    undefined
)

export const useAuth = () => {
    const ctx = useContext(AuthContext)
    if (!ctx) {
        throw new Error('useAuth must be used inside AuthProvider')
    }
    return ctx
}
