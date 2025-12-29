import React from 'react';

type CardProps = {
    title?: string,
    width?: number
    children: React.ReactNode
}

export const Card = ({ title, width, children }: CardProps) => {
    return (
        <div className={`bg-[var(--surface)]
    text-[var(--text)]
    p-[var(--space-xl)]
    border border-[var(--border)]
    shadow-[var(--shadow-md)]
    rounded-[var(--radius-lg)]`}
             style={{width:width ?? 'auto'}}>
            {title && <h2 className={'mt-0 mb-[var(--space-sm)] text-[length:var(--text-xl)] font-semibold'}>{title}</h2>}
            {children}
        </div>
    );
};
