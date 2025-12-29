import React, { useRef, useState } from 'react';
import { MessageContext } from './Context'
import type {MessageLevel, MessageState} from './types';

export function MessageProvider({ children }: { children: React.ReactNode }) {
    const [message, setMessage] = useState<MessageState>({ text: '', type: 'info' });
    const timeoutRef = useRef<number | null>(null);

    const clearMessage = () => setMessage({ text: '', type: 'info' });

    const showMessage = (text: string, type?: MessageLevel, duration?: number) => {
        if (!text) return;

        type = type || 'info';
        if (typeof duration === 'undefined') {
            duration = 15000;
        }

        setMessage({ text, type });

        if (duration > 0) {
            if (timeoutRef.current) clearTimeout(timeoutRef.current);
            timeoutRef.current = window.setTimeout(clearMessage, duration);
        }
    };

    return (
        <MessageContext.Provider value={{ message, showMessage, clearMessage }}>
            {children}
        </MessageContext.Provider>
    );
}
