import React from "react";
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { BrowserRouter } from 'react-router-dom';

import { LoaderProvider } from '../LoaderProvider';
import { MessageProvider } from '../MessageProvider';
import { ConfigProvider } from '../ConfigProvider';
import { LangProvider } from '../LangProvider';
import { AuthProvider } from '../AuthProvider';
import { AppThemeProvider } from '../AppThemeProvider';

const queryClient = new QueryClient();

type Props = {
    children: React.ReactNode;
};

export function AppProviders({ children } : Props) {
  return (
    <BrowserRouter>
      <QueryClientProvider client={queryClient}>
        <LoaderProvider>
          <MessageProvider>
            <ConfigProvider>
              <LangProvider>
                <AuthProvider>
                    <AppThemeProvider>
                        {children}
                    </AppThemeProvider>
                </AuthProvider>
              </LangProvider>
            </ConfigProvider>
          </MessageProvider>
        </LoaderProvider>
      </QueryClientProvider>
    </BrowserRouter>
  );
}
