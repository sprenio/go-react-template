import React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { BrowserRouter } from 'react-router-dom';

import { LoaderProvider } from '../LoaderProvider';
import { MessageProvider } from '../MessageProvider';
import { ConfigProvider } from '../ConfigProvider';
import { LangProvider } from '../LangProvider';
import { AuthProvider } from '../AuthProvider';
import { AppThemeProvider } from '../AppThemeProvider';
import { VersionWatcherProvider } from '../VersionWatcherProvider';
import { TooltipProvider } from '@/components/ui/tooltip';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      gcTime: 1000 * 60 * 60 * 24, // 24 hours
    },
  },
});

type Props = {
  children: React.ReactNode;
};

export function AppProviders({ children }: Props) {
  return (
    <BrowserRouter>
      <QueryClientProvider client={queryClient}>
        <LoaderProvider>
          <MessageProvider>
            <ConfigProvider>
              <LangProvider>
                <AuthProvider>
                  <TooltipProvider>
                    <AppThemeProvider>
                      <VersionWatcherProvider>{children}</VersionWatcherProvider>
                    </AppThemeProvider>
                  </TooltipProvider>
                </AuthProvider>
              </LangProvider>
            </ConfigProvider>
          </MessageProvider>
        </LoaderProvider>
      </QueryClientProvider>
    </BrowserRouter>
  );
}
