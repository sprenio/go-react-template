import React, { useEffect, useState } from 'react';
import { useIsMutating } from '@tanstack/react-query';
import { setLoaderHandlers } from '@/api';
import Loader from './Loader';
import { LoaderContext } from './Context';

export function LoaderProvider({ children }: { children: React.ReactNode }) {
  const [loader, setLoader] = useState({ visible: false, message: '' });
  const [manualMode, setManualMode] = useState(false);

  const isMutating = useIsMutating();

  const showLoader = (message = '') => {
    setManualMode(true);
    setLoader({ visible: true, message });
  };

  const hideLoader = () => {
    setManualMode(false);
    setLoader({ visible: false, message: '' });
  };
  setLoaderHandlers({ showLoader, hideLoader });

  useEffect(() => {
    if (manualMode) return;
    setLoader({
      visible: isMutating > 0,
      message: '',
    });
  }, [isMutating, manualMode]);

  return (
    <LoaderContext.Provider value={{ loader, showLoader, hideLoader }}>
      {loader.visible && <Loader message={loader.message} />}
      {children}
    </LoaderContext.Provider>
  );
}
