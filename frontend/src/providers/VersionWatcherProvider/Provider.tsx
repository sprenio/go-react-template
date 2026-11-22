import { useState, useEffect, useRef, useCallback } from 'react';
import { useLocation } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { toast } from 'react-hot-toast';
import { setVersionMismatchHandler } from '@/api/apiClient';
import { Button } from '@/components/Button';
import { VersionWatcherContext } from './Context';
import type { VersionWatcherProviderProps } from './types';

const TOAST_ID = 'app-version-update-toast';

export function VersionWatcherProvider({ children }: VersionWatcherProviderProps) {
  const { t } = useTranslation();
  const location = useLocation();
  const [isUpdateAvailable, setIsUpdateAvailable] = useState<boolean>(false);
  const isUpdateAvailableRef = useRef<boolean>(false);
  const isInitialLocationRef = useRef<boolean>(true);

  const refreshApplication = useCallback(() => {
    window.location.reload();
  }, []);

  const handleVersionMismatch = useCallback(() => {
    if (isUpdateAvailableRef.current) {
      return;
    }

    isUpdateAvailableRef.current = true;
    setIsUpdateAvailable(true);

    toast(
      () => (
        <div className="flex items-center gap-3">
          <span className="text-sm font-medium">
            {t('global.version_update_available')}
          </span>
          <Button
            type="primary"
            size="sm"
            onClick={refreshApplication}
          >
            {t('global.btn_refresh')}
          </Button>
        </div>
      ),
      {
        id: TOAST_ID,
        duration: Infinity,
      },
    );
  }, [refreshApplication, t]);

  useEffect(() => {
    setVersionMismatchHandler(handleVersionMismatch);
    return () => {
      setVersionMismatchHandler(null);
    };
  }, [handleVersionMismatch]);

  useEffect(() => {
    if (isInitialLocationRef.current) {
      isInitialLocationRef.current = false;
      return;
    }

    if (isUpdateAvailableRef.current) {
      refreshApplication();
    }
  }, [location.pathname, location.search, refreshApplication]);

  return (
    <VersionWatcherContext.Provider value={{ isUpdateAvailable, refreshApplication }}>
      {children}
    </VersionWatcherContext.Provider>
  );
}

