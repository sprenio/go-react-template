import { useEffect, useState } from 'react';
import { PageTitle } from '@/layouts/PageTitle';
import { useAuth } from '@/providers/AuthProvider';
import { useSettings, useEmailChange } from '@/hooks/useSettings';
import type { UserType, UserSettingsType, AppSettingsType } from '@/domain';
import { useLang } from '@/providers/LangProvider';
import { useTranslation } from 'react-i18next';
import { toastError, ApiError } from '@/api';
import { GridSmLgWrapper } from '@/layouts/GridSmLgWrapper';
import { Card } from '@/components/Card';
import UserSettingsForm from './UserSettingsForm';
import AppSettingsForm from './AppSettingsForm';
import EmailChangeForm from './EmailChangeForm';
import { NotFound } from '@/pages/NotFound';
import { defUserSettings, defAppSettings, defSettings } from './defaults';
import { useMessage } from '@/providers/MessageProvider';

function getUserSettings(user: UserType | null): UserSettingsType {
  return user
    ? {
      language: user.settings?.language || defUserSettings.language,
      user_flags: user.settings?.user_flags || defUserSettings.user_flags,
    }
    : defUserSettings;
}

function getAppSettings(user: UserType | null): AppSettingsType {
  return user
    ? {
      app_opt_1: user.settings?.app_opt_1 || defAppSettings.app_opt_1,
      app_opt_2: user.settings?.app_opt_2 || defAppSettings.app_opt_2,
      app_opt_3: user.settings?.app_opt_3 || defAppSettings.app_opt_3,
      app_flags: user.settings?.app_flags || defAppSettings.app_flags,
    }
    : defAppSettings;
}

export default function Settings() {
  const [activeSection, setActiveSection] = useState('profile');
  const { appUser, setLoginUser, fetchMe } = useAuth();
  const settingsMutation = useSettings();
  const emailChangeMutation = useEmailChange();
  const { showMessage } = useMessage();
  const { changeLang } = useLang();
  const [initUserSettings, setInitUserSettings] = useState(getUserSettings(appUser));

  const [initAppSettings, setInitAppSettings] = useState(getAppSettings(appUser));
  const { t } = useTranslation();

  useEffect(() => {
    fetchMe();
  }, []);
  if (!appUser) {
    return <NotFound />;
  }
  const handleSettingsSave = (settings: UserSettingsType | AppSettingsType) => {
    if (!appUser) {
      return;
    }
    const newSettings = { ...defSettings, ...appUser.settings, ...settings };
    settingsMutation.mutate(newSettings, {
      onSuccess: () => {
        showMessage(t('settings.settings_saved'), 'success');
        const langChanged = newSettings.language !== appUser.settings?.language;
        const updatedUser = { ...appUser, settings: newSettings };
        setLoginUser(updatedUser);
        if (langChanged) {
          changeLang(newSettings.language);
        }
        setInitUserSettings(getUserSettings(updatedUser));
        setInitAppSettings(getAppSettings(updatedUser));
      },
      onError: (error: ApiError) => {
        console.error('Error saving user settings:', error);
        toastError(error);
      },
    });
  };

  const handleEmailChange = (email: string) => {
    if (!email) {
      return;
    }
    emailChangeMutation.mutate(
      { email },
      {
        onSuccess: () => {
          showMessage(t('settings.confirmation_email_sent'), 'success');
        },
        onError: (error: ApiError) => {
          console.error('Error saving user email:', error);
          toastError(error);
        },
      },
    );
  };

  return (
    <>
      <PageTitle page={'settings'} />
      <GridSmLgWrapper>
        <div className="grid-column">
          <Card>
            {/* SIDEBAR */}
            <aside className={'sidebar'}>
              <nav>
                <ul>
                  <li>
                    <button
                      className={`item ${activeSection === 'profile' ? 'active' : ''}`}
                      onClick={() => setActiveSection('profile')}
                    >
                      {t('settings.profile')}
                    </button>
                  </li>
                  <li>
                    <button
                      className={`item ${activeSection === 'app' ? 'active' : ''}`}
                      onClick={() => setActiveSection('app')}
                    >
                      {t('settings.app_settings')}
                    </button>
                  </li>
                  <li>
                    <button
                      className={`item ${activeSection === 'email' ? 'active' : ''}`}
                      onClick={() => setActiveSection('email')}
                    >
                      {t('settings.email')}
                    </button>
                  </li>
                </ul>
              </nav>
            </aside>
          </Card>
        </div>
        <div className="grid-column">
          {/* MAIN CONTENT */}
          <div className={'sidebar-content'}>
            {activeSection === 'profile' && (
              <Card title={t('settings.profile')}>
                <UserSettingsForm
                  userSettings={initUserSettings}
                  user={appUser}
                  onSubmit={handleSettingsSave}
                />
              </Card>
            )}

            {activeSection === 'app' && (
              <Card title={t('settings.app_settings')}>
                <AppSettingsForm appSettings={initAppSettings} onSubmit={handleSettingsSave} />
              </Card>
            )}

            {activeSection === 'email' && (
              <Card title={t('settings.email_change')}>
                <EmailChangeForm onSubmit={handleEmailChange} email={appUser?.email} />
              </Card>
            )}
          </div>
        </div>
      </GridSmLgWrapper>
    </>
  );
}
