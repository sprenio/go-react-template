import React, { useEffect, useState } from 'react';
import { api , toastError, ApiError } from '@/api';
import {ConfigContext} from './Context';
import type {ConfigType} from '@/types'

export function ConfigProvider({ children }: { children: React.ReactNode }) {
    const [loading, setLoading] = useState<boolean>(true);
    const [config, setConfig] = useState<ConfigType>({
        Features: {},
        AppName: import.meta.env.VITE_APP_NAME || 'My Application',
        Languages: [],
        DefaultLanguage: 'pl'
    });
    const [showLanguages, setShowLanguages] = useState<boolean>(false)

    useEffect(() => {
        api
            .get<ConfigType>('/cfg')
            .then(({data}) => {
                setConfig(data);
                setShowLanguages(data.Languages.length > 1)
            })
            .catch((err:ApiError) => {
                console.error('❌ Błąd pobierania konfiguracji:', err);
                toastError(err);
            })
            .finally(() => setLoading(false));
    }, []);

    return <ConfigContext.Provider value={{ config, loading, showLanguages }}>{children}</ConfigContext.Provider>;
}
