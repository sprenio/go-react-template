export type MessageLevel = 'info' | 'success' | 'warning' | 'error';

export type MessageState = {
    text: string
    type: MessageLevel
    duration?: number
}

export type MessageContextType = {
    message: MessageState
    showMessage: (text: string, type:MessageLevel, duration?: number) => void
    clearMessage: () => void
}
