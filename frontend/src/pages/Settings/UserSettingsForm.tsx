import {LanguageSelect} from '@/components/LanguageSelect';
import {FormWrapper} from '@/components/Form';
import {useTranslation} from 'react-i18next';
import React, {useReducer} from 'react';
import type {UserType, UserSettingsType} from '@/types';
import {defUserSettings} from './defaults';
import { useConfig } from '@/providers/ConfigProvider';

const reducerActions = {
    SET_LANGUAGE: 'SET_LANGUAGE',
    SET_FLAGS: 'SET_FLAGS',
};

type UserSettingsFormProps = {
    userSettings: UserSettingsType
    user: UserType
    onSubmit: (settings: UserSettingsType) => void
}
type SettingsReducerAction = {
    type: string,
    flag: string,
    value: string,
}

const settingsReducer = (state: UserSettingsType, action: SettingsReducerAction) => {
    switch (action.type) {
        case reducerActions.SET_LANGUAGE:
            return {...state, language: action.value};
        case reducerActions.SET_FLAGS:
            state.user_flags[action.flag] = !!action.value;
            return {
                ...state,
                user_flags: {
                    ...state.user_flags,
                    [action.flag]: !!action.value,
                },
            };
        default:
            return state;
    }
};

export default function UserSettingsForm({userSettings, user, onSubmit}: UserSettingsFormProps) {
    const [settings, dispatch] = useReducer(
        settingsReducer,
        userSettings || {
            language: user.settings?.language || defUserSettings.language,
            user_flags: user.settings?.user_flags || defUserSettings.user_flags,
        }
    );

    const {t} = useTranslation();

    const handleSaveSettings = (e:React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        if (onSubmit) {
            onSubmit(settings);
        }
    };

    const handleFlagChange = (e:React.ChangeEvent<HTMLInputElement>) => {
        dispatch({
            type: reducerActions.SET_FLAGS,
            flag: e.target.id,
            value: e.target.checked ? '1' : '',
        });
    };
    const handleLanguageChange = (lang: string) => {
        dispatch({
            type: reducerActions.SET_LANGUAGE,
            value: lang,
            flag: '',
        });
    };

    const { showLanguages } = useConfig();
    return (
        <FormWrapper onSubmit={handleSaveSettings}>
            {/* User name */}
            <div>
                <label>{t('settings.lbl_user_name')}</label>
                <p className={'py-2 text-[length:var(--text-lg)]'}>{user.name}</p>
            </div>

            {/* Email */}
            <div>
                <label>{t('settings.lbl_email')}</label>
                <p className={'py-2 text-[length:var(--text-lg)]'}>{user.email}</p>
            </div>

            {/* Language */}
            {showLanguages &&
            <div>
                <label>{t('settings.lbl_language')}</label>
                <LanguageSelect value={settings.language} onChange={handleLanguageChange}/>
            </div>
            }

            {/* Flags */}
            <div>
                <label>{t('settings.lbl_extra_options')}</label>
                <div>
                    <label className="checkbox-label">
                        <input
                            type="checkbox"
                            id="flag_1"
                            checked={settings.user_flags.flag_1}
                            onChange={handleFlagChange}
                        />
                        {t('settings.flag_1')}
                    </label>
                </div>
                <div>
                    <label className="checkbox-label">
                        <input
                            type="checkbox"
                            id="flag_2"
                            checked={settings.user_flags.flag_2}
                            onChange={handleFlagChange}
                        />
                        {t('settings.flag_2')}
                    </label>
                </div>
                <div>
                    <label className="checkbox-label">
                        <input
                            type="checkbox"
                            id="flag_3"
                            checked={settings.user_flags.flag_3}
                            onChange={handleFlagChange}
                        />
                        {t('settings.flag_3')}
                    </label>
                </div>
            </div>

            <button type="submit">{t('settings.btn_save')}</button>
        </FormWrapper>
    );
}
