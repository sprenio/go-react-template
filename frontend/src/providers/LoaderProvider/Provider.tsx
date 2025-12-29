import React, { useEffect, useState } from 'react';
import { useIsFetching, useIsMutating } from '@tanstack/react-query';
import { setLoaderHandlers } from '@/api';
import Loader from './Loader';
import { LoaderContext } from './Context';

export function LoaderProvider({ children }: { children: React.ReactNode }) {
    const [loader, setLoader] = useState({ visible: false, message: '' })
    const [manualMode, setManualMode] = useState(false)

    const isFetching = useIsFetching()
    const isMutating = useIsMutating()

    const showLoader = (message = '') => {
        setManualMode(true)
        setLoader({ visible: true, message })
    }

    const hideLoader = () => {
        setManualMode(false)
        setLoader({ visible: false, message: '' })
    }
    useEffect(() => {
        setLoaderHandlers({ showLoader, hideLoader });
    }, []);

    useEffect(() => {
        if (manualMode) return
        setLoader({
            visible: isFetching > 0 || isMutating > 0,
            message: '',
        })
    }, [isFetching, isMutating, manualMode])

    return (
        <LoaderContext.Provider value={{ loader, showLoader, hideLoader }}>
            {loader.visible && <Loader message={loader.message} />}
            {children}
        </LoaderContext.Provider>
    )
}
