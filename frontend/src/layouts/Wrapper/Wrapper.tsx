import React from 'react';

export function Wrapper({ children }: { children: React.ReactNode }) {
    return (
        <div className="flex min-h-screen flex-col overflow-hidden bg-[var(--bg)]">
            {children}
        </div>
    )
}
