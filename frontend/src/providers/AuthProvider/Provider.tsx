import React, {useState, useEffect, useRef} from 'react';
import {useNavigate} from 'react-router-dom';
import {useTranslation} from 'react-i18next';
import {useLang} from '../LangProvider'
import {AuthContext} from './Context'
import {api, setLogoutHandler, logoutReasonCodes, type ApiSuccessResponse} from '@/api';
import type {MeResponse} from './types';
import type {UserType} from '@/types';
import {setSuccessHandler} from "@/api/apiClient.ts";

export function AuthProvider({children}: { children: React.ReactNode }) {
    const {t} = useTranslation();
    const [appUser, setAppUser] = useState<UserType | null>(null);
    const navigate = useNavigate();
    const [meInProgress, setMeInProgress] = useState<boolean>(true);
    const refreshTimeout = useRef<number | null>(null);

    const {changeLang} = useLang();

    const appUserRef = useRef<UserType | null>(null);

    useEffect(() => {
        appUserRef.current = appUser;
    }, [appUser]);

    useEffect(() => {
        setLogoutHandler(logout);
        setSuccessHandler(apiSuccessResponseHandler)
        fetchMe().finally(() => setMeInProgress(false));
        return () => {
            if (refreshTimeout.current) clearTimeout(refreshTimeout.current);
        };
    }, []);

    const fetchMe = async () => {
        try {
            await api.get<MeResponse>('/me', false).then(({data}) => {
                setLoginUser(data.user);
            });
        } catch (err) {
            console.error('Refresh /me failed', err);
            if(appUserRef.current) {
                logout(logoutReasonCodes.SESSION_EXPIRED);
            }
        }
    };
    const clearRefreshTimeout = () => {
        if (refreshTimeout.current) {
            clearTimeout(refreshTimeout.current);
        }
    }
    const apiSuccessResponseHandler = () => {
        if (appUserRef.current) {
            scheduleMeRefresh();
        }
    }
    const scheduleMeRefresh = () => {

        // Clear previous timeout if exists
        clearRefreshTimeout()
        // set next call /me request after 5 minutes (300_000 ms)
        refreshTimeout.current = window.setTimeout(async () => {
            try {
                await fetchMe();
            } finally {
                if(appUserRef.current) {
                    scheduleMeRefresh();
                }
            }
        }, 30_000);
    };

    const setLoginUser = (user: UserType | null) => {
        setAppUser(user);
        const lang = user?.settings?.language;
        if (lang) {
            changeLang(lang);
        }
        if (user) {
            scheduleMeRefresh();
        } else {
            clearRefreshTimeout()
        }
    };

    const logout = (REASON_CODE: string) => {
        if(!appUserRef.current){
            return
        }
        clearRefreshTimeout();
        setAppUser(null);
        let state = {};

        switch (REASON_CODE) {
            case logoutReasonCodes.SESSION_EXPIRED:
                state = {message: t('global.session_expired'), type: 'error'};
                break;
            case logoutReasonCodes.USER_LOGGED_OUT:
                state = {message: t('global.logout_success'), type: 'success'};
                break;
            default:
                state = {message: t('global.unknown_error'), type: 'error'};
                break;
        }
        try {
            api.get<ApiSuccessResponse>('/logout').then(() => {
                navigate('/login', state);
            });
        } catch (err) {
            console.error('Logout failed', err);
        }
    };

    return (
        <AuthContext.Provider
            value={{
                setLoginUser,
                logout,
                appUser,
                meInProgress,
            }}
        >
            {children}
        </AuthContext.Provider>
    );
}