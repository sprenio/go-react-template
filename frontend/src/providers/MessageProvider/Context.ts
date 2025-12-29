import { createContext, useContext } from 'react'
import type { MessageContextType } from './types'

export const MessageContext = createContext<MessageContextType | undefined>(
    undefined
)

export const useMessage = () => {
    const ctx = useContext(MessageContext)
    if (!ctx) {
        throw new Error('useMessage must be used inside MessageProvider')
    }
    return ctx
}
