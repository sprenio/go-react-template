import { createContext, useContext } from 'react';
import type { VersionWatcherContextType } from './types';

export const VersionWatcherContext = createContext<VersionWatcherContextType | undefined>(undefined);

export function useVersionWatcher(): VersionWatcherContextType {
  const context = useContext(VersionWatcherContext);
  if (!context) {
    throw new Error('useVersionWatcher must be used inside VersionWatcherProvider');
  }
  return context;
}

