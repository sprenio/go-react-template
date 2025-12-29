import React, {useReducer} from 'react';
import {FormWrapper} from '@/components/Form';
import {useTranslation} from 'react-i18next';
import type {AppSettingsType} from '@/types';

const reducerActions = {
    SET_FLAGS: 'SET_FLAGS',
    SET_OPTION: 'SET_OPTION',
};

type AppSettingsFormProps = {
    appSettings: AppSettingsType
    onSubmit: (settings: AppSettingsType) => void
}
type SettingsReducerAction = {
    type: string,
    flag: string,
    option: string,
    value: string,
}

const settingsReducer = (state: AppSettingsType, action: SettingsReducerAction) => {
    switch (action.type) {
        case reducerActions.SET_OPTION:
            return {...state, [action.option]: action.value};
        case reducerActions.SET_FLAGS:
            state.app_flags[action.flag] = !!action.value;
            return {
                ...state,
                app_flags: {
                    ...state.app_flags,
                    [action.flag]: !!action.value,
                },
            };
        default:
            return state;
    }
};

export default function AppSettingsForm({appSettings, onSubmit}: AppSettingsFormProps) {
    const [settings, dispatch] = useReducer(settingsReducer, appSettings);

    const {t} = useTranslation();
    const handleSaveSettings = (e:React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        if (onSubmit) {
            onSubmit(settings);
        }
    };

    const handleOptionChange = (e:React.ChangeEvent<HTMLSelectElement|HTMLInputElement>) => {
        dispatch({
            type: reducerActions.SET_OPTION,
            option: e.target.name,
            value: e.target.value,
            flag: '',
        });
    };
    const handleFlagChange = (e:React.ChangeEvent<HTMLInputElement>) => {
        dispatch({
            type: reducerActions.SET_FLAGS,
            flag: e.target.name,
            value: e.target.checked ? '1' : '',
            option: '',
        });
    };

    return (
        <FormWrapper onSubmit={handleSaveSettings}>
            {/* input:text */}
            <div>
                <label>{t('settings.lbl_app_opt_1')}</label>
                <input
                    type="text"
                    placeholder={t('settings.ph_app_opt_1')}
                    name="app_opt_1"
                    value={settings.app_opt_1}
                    onChange={handleOptionChange}
                />
            </div>

            {/* select */}
            <div>
                <label>{t('settings.lbl_app_opt_2')}</label>
                <div className="select-wrapper">
                    <select name="app_opt_2" value={settings.app_opt_2} onChange={handleOptionChange}>
                        <option value="">{t('settings.option_none')}</option>
                        <option value="OPT_A">{t('settings.option_a')}</option>
                        <option value="OPT_B">{t('settings.option_b')}</option>
                    </select>
                </div>
            </div>

            {/* checkbox */}
            <div>
                <label>{t('settings.lbl_app_opt_3')}</label>
                <div>
                    <label className="checkbox-label">
                        <input
                            type="checkbox"
                            name="flag_a"
                            checked={settings.app_flags.flag_a}
                            onChange={handleFlagChange}
                        />
                        {t('settings.flag_a')}
                    </label>
                </div>
                <div>
                    <label className="checkbox-label">
                        <input
                            type="checkbox"
                            name="flag_b"
                            checked={settings.app_flags.flag_b}
                            onChange={handleFlagChange}
                        />
                        {t('settings.flag_b')}
                    </label>
                </div>
            </div>

            {/* radio */}
            <div>
                <label>{t('settings.lbl_app_opt_4')}</label>
                <div>
                    <label className="radio-label">
                        <input
                            type="radio"
                            name="app_opt_3"
                            value="RADIO_A"
                            checked={settings.app_opt_3 === 'RADIO_A'}
                            onChange={handleOptionChange}
                        />
                        {t('settings.option_1')}
                    </label>
                </div>
                <div>
                    <label className="radio-label">
                        <input
                            type="radio"
                            name="app_opt_3"
                            value="RADIO_B"
                            checked={settings.app_opt_3 === 'RADIO_B'}
                            onChange={handleOptionChange}
                        />
                        {t('settings.option_2')}
                    </label>
                </div>
            </div>

            <button type="submit">{t('settings.btn_save')}</button>
        </FormWrapper>
    );
}
