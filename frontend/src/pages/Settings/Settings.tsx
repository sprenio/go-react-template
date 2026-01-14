import {useState} from 'react';
import {PageTitle} from '@/layouts/PageTitle';
import {useAuth} from '@/providers/AuthProvider';
import {useSettings, useEmailChange} from '@/hooks/useSettings';
import type {UserType, UserSettingsType, AppSettingsType} from '@/types';
import {useLang} from '@/providers/LangProvider';
import {useTranslation} from 'react-i18next';
import {toastError, toastSuccess} from '@/api';
import {GridSmLgWrapper} from '@/layouts/GridSmLgWrapper';
import {Card} from '@/components/Card';
import UserSettingsForm from './UserSettingsForm';
import AppSettingsForm from './AppSettingsForm';
import EmailChangeForm from './EmailChangeForm';
import {NotFound} from '@/pages/NotFound';
import {defUserSettings, defAppSettings, defSettings} from './defaults';

function getUserSettings(user: UserType | null): UserSettingsType {
    return user ? {
        language: user.settings?.language || defUserSettings.language,
        user_flags: user.settings?.user_flags || defUserSettings.user_flags,
    } : defUserSettings;
}

function getAppSettings(user: UserType | null): AppSettingsType {
    return user ? {
        app_opt_1: user.settings?.app_opt_1 || defAppSettings.app_opt_1,
        app_opt_2: user.settings?.app_opt_2 || defAppSettings.app_opt_2,
        app_opt_3: user.settings?.app_opt_3 || defAppSettings.app_opt_3,
        app_flags: user.settings?.app_flags || defAppSettings.app_flags,
    } : defAppSettings;
}

export default function Settings() {
    const [activeSection, setActiveSection] = useState('profile');
    const {appUser, setLoginUser} = useAuth();
    const settingsMutation = useSettings();
    const emailChangeMutation = useEmailChange();

    const {changeLang} = useLang();
    const [initUserSettings, setInitUserSettings] = useState(getUserSettings(appUser));

    const [initAppSettings, setInitAppSettings] = useState(getAppSettings(appUser));
    if (!appUser) {
        return (<NotFound/>)
    }
    const handleSettingsSave = (settings: UserSettingsType | AppSettingsType) => {
        if (!appUser) {
            return
        }
        const newSettings = {...defSettings, ...appUser.settings, ...settings};
        settingsMutation.mutate(
            newSettings,
            {
                onSuccess: (resp) => {
                    toastSuccess(resp.code);
                    const langChanged = newSettings.language !== appUser.settings?.language;
                    const updatedUser = {...appUser, settings: newSettings};
                    setLoginUser(updatedUser);
                    if (langChanged) {
                        changeLang(newSettings.language);
                    }
                    setInitUserSettings(getUserSettings(updatedUser));
                    setInitAppSettings(getAppSettings(updatedUser));
                },
                onError: (error) => {
                    console.error('Error saving user settings:', error);
                    toastError(error);
                }
            }
        );
    };

    const handleEmailChange = (email: string) => {
        if (!email) {
            return;
        }
        emailChangeMutation.mutate(
            {email},
            {
                onSuccess: (resp) => {
                    toastSuccess(resp.code);
                },
                onError: (error) => {
                    console.error('Error saving user email:', error);
                    toastError(error);
                }
            }
        );
    };

    const {t} = useTranslation();

    return (
        <>
            <PageTitle page={'settings'}/>
            <GridSmLgWrapper>
                <div className="grid-column">
                    <Card>
                        {/* SIDEBAR */}
                        <aside className={'settings-sidebar'}>
                            <nav>
                                <ul>
                                    <li>
                                        <button
                                            className={activeSection === 'profile' ? 'active' : ''}
                                            onClick={() => setActiveSection('profile')}
                                        >
                                            {t('settings.profile')}
                                        </button>
                                    </li>
                                    <li>
                                        <button
                                            className={activeSection === 'app' ? 'active' : ''}
                                            onClick={() => setActiveSection('app')}
                                        >
                                            {t('settings.app_settings')}
                                        </button>
                                    </li>
                                    <li>
                                        <button
                                            className={activeSection === 'email' ? 'active' : ''}
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
                    <div className={'settings-content'}>
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
                                <AppSettingsForm appSettings={initAppSettings} onSubmit={handleSettingsSave}/>
                            </Card>
                        )}

                        {activeSection === 'email' && (
                            <Card title={t('settings.email_change')}>
                                <EmailChangeForm onSubmit={handleEmailChange} email={appUser?.email}/>
                            </Card>
                        )}
                    </div>
                </div>
            </GridSmLgWrapper>
        </>
    );
}
