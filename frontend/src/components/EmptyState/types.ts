import type { ReactNode } from 'react';

export type EmptyStateType = 'info' | 'warning' | 'error';

export type EmptyStateProps = {
  message?: string;
  type?: EmptyStateType;
  action?: ReactNode;
};
