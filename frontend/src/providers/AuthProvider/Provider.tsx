import React, {useState, useEffect, useRef} from 'react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useLang } from '../LangProvider'
import { AuthContext } from './Context'
import { api, setLogoutHandler, setTokenRefresher, logoutReasonCodes } from '@/api';
import {
    clearToken,
    getRememberMeFlag,
    getToken,
    saveToken,
    isLoggedIn,
} from './auth';
import type {MeResponse} from './types';
import type {UserType} from '@/types';
export function AuthProvider({ children }: { children: React.ReactNode }) {
    const { t } = useTranslation();
    const [token, setToken] = useState<string | null>(null);
    const [user, setUser] = useState<UserType | null>(null);
    const navigate = useNavigate();
    const [meInProgress, setMeInProgress] = useState<boolean>(true);
    const refreshTimeout = useRef<number | null>(null);

    const { changeLang } = useLang();

    useEffect(() => {
        setLogoutHandler(logout);
        setTokenRefresher(setLoginToken);

        // Ładuj z localStorage przy starcie aplikacji
        const storedToken = getToken();
        if (storedToken) {
            fetchMe().finally(() => setMeInProgress(false));
        } else {
            setMeInProgress(false);
        }
        return () => {
            if (refreshTimeout.current) clearTimeout(refreshTimeout.current);
        };
    }, []);

    const fetchMe = async () => {
        try {
            const resp = await api.get<MeResponse>('/me');
            if (resp?.token) setLoginToken(resp.token);
            if (resp?.data?.user) setUser(resp.data.user);
        } catch (err) {
            console.error('Refresh /me failed', err);
            logout(logoutReasonCodes.SESSION_EXPIRED);
        }
    };
    const scheduleMeRefresh = () => {
        // Czyścimy poprzedni timeout (jeśli istnieje)
        if (refreshTimeout.current) {
            clearTimeout(refreshTimeout.current);
        }
        if (isLoggedIn()) {
            // ustawiamy kolejne wywołanie /me po 5 minutach (300_000 ms)
            refreshTimeout.current = window.setTimeout(async () => {
                try {
                    await fetchMe();
                } finally {
                    // ustawiamy kolejny timeout tylko jeśli jeszcze nie został zresetowany przez setLoginToken
                    scheduleMeRefresh();
                }
            }, 300_000);
        }
    };

    const setLoginUser = (user:UserType|null) => {
        setUser(user);
        const lang = user?.settings?.language;
        if (lang) {
            changeLang(lang);
        }
    };
    const setLoginToken = (token:string) => {
        setToken(token);
        const rememberMe = getRememberMeFlag();
        saveToken(token, rememberMe);
        // reset timer po każdym odświeżeniu tokena
        if (refreshTimeout.current) clearTimeout(refreshTimeout.current);
        scheduleMeRefresh();
    };

    const logout = (REASON_CODE:string) => {
        setToken(null);
        setUser(null);
        clearToken();
        let state = {};

        switch (REASON_CODE) {
            case logoutReasonCodes.SESSION_EXPIRED:
                state = { message: t('global.session_expired'), type: 'error' };
                break;
            case logoutReasonCodes.USER_LOGGED_OUT:
                state = { message: t('global.logout_success'), type: 'success' };
                break;
            default:
                state = { message: t('global.unknown_error'), type: 'error' };
                break;
        }
        navigate('/login', state);
    };

    return (
        <AuthContext.Provider
            value={{
                token,
                setLoginUser,
                setLoginToken,
                logout,
                user,
                meInProgress,
            }}
        >
            {children}
        </AuthContext.Provider>
    );
}