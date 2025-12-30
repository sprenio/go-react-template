import {render} from '@testing-library/react';
import {QueryClient, QueryClientProvider} from '@tanstack/react-query';
import {AuthContext} from '@/providers/AuthProvider';
import {ConfigContext} from '@/providers/ConfigProvider';
import {LoaderProvider} from '@/providers/LoaderProvider';
import {LangProvider} from '@/providers/LangProvider';
import {MessageProvider} from '@/providers/MessageProvider';
import {AppThemeProvider} from '@/providers/AppThemeProvider';
import {MemoryRouter} from 'react-router-dom';
import React from 'react';
import type {JSX} from 'react';

export function renderWithProviders(
    ui: JSX.Element,
    {
        route = '/',
        auth = {
            token: '',
            setLoginUser: vi.fn(),
            setLoginToken: vi.fn(),
            logout: vi.fn(),
            user: null,
            meInProgress: false
        },
        featureFlags = {flags: {register: true, reset_password: true}},
    } = {}
) {
    const queryClient = new QueryClient();

    // Mock ConfigProvider with the provided feature flags
    const MockConfigProvider = ({children}: { children: React.ReactNode }) => {
        const mockConfig = {
            Features: featureFlags.flags || {register: true, reset_password: true},
            AppName: 'Test App',
            Languages: [],
            DefaultLanguage: 'pl'
        };
        return (
            <ConfigContext.Provider value={{config: mockConfig, loading: false}}>
                {children}
            </ConfigContext.Provider>
        );
    };

    return render(
        <MemoryRouter initialEntries={[route]}>
            <QueryClientProvider client={queryClient}>
                <LoaderProvider>
                    <MessageProvider>
                        <MockConfigProvider>
                            <LangProvider>
                                <AuthContext.Provider value={auth}>
                                    <AppThemeProvider>
                                        <AppThemeProvider>{ui}</AppThemeProvider>
                                    </AppThemeProvider>
                                </AuthContext.Provider>
                            </LangProvider>
                        </MockConfigProvider>
                    </MessageProvider>
                </LoaderProvider>
            </QueryClientProvider>
        </MemoryRouter>
    );
}
