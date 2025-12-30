import {renderWithProviders} from '@/test/renderWithProviders';
import App from '../App';
import {screen} from '@testing-library/react';

// Mock localStorage and sessionStorage
const localStorageMock = {
    getItem: vi.fn(),
    setItem: vi.fn(),
    removeItem: vi.fn(),
    clear: vi.fn(),
};
const sessionStorageMock = {
    getItem: vi.fn(),
    setItem: vi.fn(),
    removeItem: vi.fn(),
    clear: vi.fn(),
};

Object.defineProperty(window, 'localStorage', {
    value: localStorageMock,
});
Object.defineProperty(window, 'sessionStorage', {
    value: sessionStorageMock,
});

// Reset mocks before each test
beforeEach(() => {
    localStorageMock.getItem.mockClear();
    sessionStorageMock.getItem.mockClear();
});

vi.mock('@/api', () => ({
    api: {
        get: vi.fn((url) => {
            if (url === '/cfg') {
                return Promise.resolve({
                    data: {
                        Features: {
                            register: true,
                            reset_password: true,
                        },
                        AppName: 'Test App',
                    },
                });
            }
            // inne endpointy jeśli chcesz, np.:
            // if (url === '/ping') { ... }
            return Promise.reject(new Error('Unknown API call'));
        }),
    },
    setLoaderHandlers: vi.fn(),
}));

const testCases = [
    {RegisterLinkFlag: true, ResetPasswordLinkFlag: true},
    {RegisterLinkFlag: true, ResetPasswordLinkFlag: false},
    {RegisterLinkFlag: false, ResetPasswordLinkFlag: true},
    {RegisterLinkFlag: false, ResetPasswordLinkFlag: false},
];
describe('When not authenticated', () => {
    test.each(testCases)(
        'shows login form with register link: $RegisterLinkFlag, reset password link: $ResetPasswordLinkFlag',
        ({RegisterLinkFlag, ResetPasswordLinkFlag}) => {
            renderWithProviders(<App/>, {
                auth: {
                    token: '',
                    setLoginUser: vi.fn(),
                    setLoginToken: vi.fn(),
                    logout: vi.fn(),
                    user: null,
                    meInProgress: false
                },
                featureFlags: {
                    flags: {
                        register: RegisterLinkFlag,
                        reset_password: ResetPasswordLinkFlag,
                    },
                },
            });

            expect(screen.getByRole('button', {name: /login.submit/})).toBeInTheDocument();

            const nav = screen.queryByRole('navigation');
            expect(nav).not.toBeInTheDocument();

            if (RegisterLinkFlag) {
                expect(screen.getByRole('link', {name: /login.register/})).toBeInTheDocument();
            } else {
                expect(screen.queryByRole('link', {name: /login.register/})).not.toBeInTheDocument();
            }

            if (ResetPasswordLinkFlag) {
                expect(screen.getByRole('link', {name: /login.forgot_password/})).toBeInTheDocument();
            } else {
                expect(
                    screen.queryByRole('link', {name: /login.forgot_password/})
                ).not.toBeInTheDocument();
            }
        }
    );
});

test('shows dashboard when authenticated', async () => {
    // Mock localStorage to return a token for isLoggedIn() function
    localStorageMock.getItem.mockReturnValue('mock-token');

    vi.stubGlobal('fetch', vi.fn(async (input: RequestInfo) => {
        if (input === '/api/ping') {
            return new Response(
                JSON.stringify({message: 'pong'}),
                {
                    status: 200,
                    headers: {'Content-Type': 'application/json'},
                }
            );
        }

        throw new Error('Unknown API call');
    }));


    renderWithProviders(<App/>, {
        route: '/dashboard',
        auth: {
            token: 'mock-token',
            setLoginUser: vi.fn(),
            setLoginToken: vi.fn(),
            logout: vi.fn(),
            user: null,
            meInProgress: false
        },
    });

    const nav = await screen.findByRole('navigation');
    expect(nav).toBeInTheDocument();
    expect(screen.queryByRole('button', {name: /login.submit/})).not.toBeInTheDocument();
});
