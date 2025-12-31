import React, {useState} from 'react';
import {useTranslation} from 'react-i18next';
import {useRegister} from '@/hooks/useRegister';
import {useNavigate} from 'react-router-dom';
import {apiCodes, getApiCodeDescription} from '@/api';
import {toast} from 'react-hot-toast';
import {LanguageSelect} from '@/components/LanguageSelect';
import {useMessage} from '@/providers/MessageProvider';
import {FormWrapper} from '@/components/Form';
import {useLang} from '@/providers/LangProvider';
import {useConfig} from '@/providers/ConfigProvider';

type RegisterFormProps = {
    onSuccess?: () => void
}

function RegisterForm({onSuccess}: RegisterFormProps) {
    const {t} = useTranslation();
    const [username, setUsername] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const registerMutation = useRegister();
    const {lang: defaultLang} = useLang();
    const [language, setLanguage] = useState(defaultLang);
    const navigate = useNavigate();
    const {showMessage} = useMessage();
    const {showLanguages} = useConfig();

    const PASSWORD_REGEX = /^(?=.*[A-Za-z])(?=.*\d)(?=.*[^A-Za-z\d]).{6,}$/;
    const minUsernameLength = 4;

    const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        if (username.length < minUsernameLength) {
            showMessage(t('register.error_username_length', {count: minUsernameLength}), 'error');
            return;
        }
        if (password !== confirmPassword) {
            showMessage(t('register.error_match'), 'error');
            return;
        }
        if (!PASSWORD_REGEX.test(password)) {
            showMessage(t('register.error_weak'), 'error');
            return;
        }

        registerMutation.mutate(
            {username, email, password, language},
            {
                onSuccess: () => {
                    setUsername('');
                    setEmail('');
                    setPassword('');
                    setConfirmPassword('');
                    onSuccess?.();
                    showMessage(t('register.success'), 'success');
                    navigate('/login');
                },
                onError: (error) => {
                    let errorMessage = getApiCodeDescription(error.code);
                    if (error.code === apiCodes.API_Register_User_Name_Or_Email_Taken) {
                        showMessage(errorMessage, 'error');
                    } else {
                        toast.error(errorMessage);
                    }
                },
            }
        );
    };

    return (
        <FormWrapper onSubmit={handleSubmit}>
            <div>
                <label htmlFor="username">{t('register.lbl_username')}</label>
                <input
                    type="text"
                    id="username"
                    autoComplete="off"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    required
                />
                <p className="helper-text">{t('register.username_hint')}</p>
            </div>

            <div>
                <label htmlFor="email">{t('register.lbl_email')}</label>
                <input
                    type="email"
                    id="email"
                    value={email}
                    autoComplete="off"
                    onChange={(e) => setEmail(e.target.value)}
                    required
                />
            </div>

            {showLanguages &&
            <div>
                <label htmlFor="username">{t('register.lbl_language')}</label>
                <LanguageSelect value={language} onChange={(lang) => setLanguage(lang)}/>
            </div>
            }

            <div>
                <label htmlFor="password">{t('register.lbl_password')}</label>
                <input
                    type="password"
                    id="password"
                    autoComplete="off"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                />
                <p className="helper-text">{t('register.password_hint')}</p>
            </div>
            <div>
                <label htmlFor="confirm">{t('register.lbl_confirm')}</label>
                <input
                    type="password"
                    id="confirm"
                    autoComplete="off"
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    required
                />
            </div>
            <button type="submit" disabled={registerMutation.isPending}>
                {registerMutation.isPending ? t('register.pending') : t('register.submit')}
            </button>
        </FormWrapper>
    );
}

export default RegisterForm;
