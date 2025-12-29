import { createContext, useContext } from 'react'
import type { LoaderContextType } from './types'

export const LoaderContext = createContext<LoaderContextType | undefined>(
    undefined
)

export const useLoader = () => {
    const ctx = useContext(LoaderContext)
    if (!ctx) {
        throw new Error('useLoader must be used inside LoaderProvider')
    }
    return ctx
}
