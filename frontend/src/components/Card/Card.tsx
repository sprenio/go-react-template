import React from 'react';

type CardProps = {
  title?: string;
  width?: number;
  titleSize?: 'xs' | 'sm' | 'md' | 'lg' | 'xl';
  children: React.ReactNode;
};

export const Card = ({ title, width, titleSize = 'xl', children }: CardProps) => {
  const titleSizeClass = `text-[length:var(--text-${titleSize})]`;
  return (
    <div
      className={`rounded-[var(--radius-lg)] border border-[var(--border)] bg-[var(--surface)] p-[var(--space-xl)] text-[var(--text)] shadow-[var(--shadow-md)]`}
      style={{ width: width ?? 'auto' }}
    >
      {title && (
        <h2 className={'mt-0 mb-[var(--space-sm)] ' + titleSizeClass + ' font-semibold'}>
          {title}
        </h2>
      )}
      {children}
    </div>
  );
};
