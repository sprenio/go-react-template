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
import type { AuthContextType } from '@/providers/AuthProvider';

type RenderOptions = {
    route?: string;
    auth?: Partial<AuthContextType>;
    featureFlags?: {
        flags?: Record<string, boolean>;
    };
};

export function renderWithProviders(
    ui: JSX.Element,
    {
        route = '/',
        auth = {},
        featureFlags = {flags: {register: true, reset_password: true}},
    }: RenderOptions = {}
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
            <ConfigContext.Provider value={{config: mockConfig, loading: false, showLanguages: true}}>
                {children}
            </ConfigContext.Provider>
        );
    };

    const authValue: AuthContextType = {
        setLoginUser: vi.fn(),
        logout: vi.fn(),
        appUser: null,
        meInProgress: false,
        ...auth, // override
    };
    return render(
        <MemoryRouter initialEntries={[route]}>
            <QueryClientProvider client={queryClient}>
                <LoaderProvider>
                    <MessageProvider>
                        <MockConfigProvider>
                            <LangProvider>
                                <AuthContext.Provider value={authValue}>
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
