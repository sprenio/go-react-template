import type { ReactNode } from 'react';

export type VersionWatcherContextType = {
  isUpdateAvailable: boolean;
  refreshApplication: () => void;
};

export type VersionWatcherProviderProps = {
  children: ReactNode;
};

