import React, { useState } from 'react';
import { useLogin } from '@/hooks/useLogin';
import { useAuth } from '@/providers/AuthProvider';
import { useTranslation } from 'react-i18next';
import { getApiCodeDescription, apiCodes } from '@/api/apiCodes';
import { toast } from 'react-hot-toast';
import { useMessage } from '@/providers/MessageProvider';
import { FormWrapper } from '@/components/Form';
import {ApiError} from '@/api';

export default function LoginForm() {
    const { t } = useTranslation();
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const loginMutation = useLogin();
    const { setLoginUser } = useAuth();
    const [rememberMe, setRememberMe] = useState(false);
    const { showMessage, clearMessage } = useMessage();

    const handleSubmit = (e:React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        clearMessage();
        loginMutation.mutate(
            { email, password, remember_me: rememberMe },
            {
                onSuccess: (data) => {
                    toast.success(t('login.success'))
                    setLoginUser(data.user);
                    setEmail('');
                },
                onError: (error:ApiError) => {
                    console.error('Login failed code:', error.code);
                    let errorMessage = getApiCodeDescription(error.code);
                    if (error.code === apiCodes.API_Login_Invalid_Credentials) {
                        showMessage(errorMessage, 'error');
                    } else {
                        toast.error(errorMessage);
                    }
                    console.error('Login failed message:', errorMessage, 'code:', error.code);
                },
                onSettled: () => {
                    setPassword('');
                },
            }
        );
    };

    return (
        <FormWrapper onSubmit={handleSubmit}>
            <div>
                <label htmlFor="email">{t('login.lbl_email')}</label>
                <input
                    id="email"
                    type="email"
                    autoComplete="off"
                    placeholder={t('login.ph_email')}
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    required
                />
            </div>
            <div>
                <label htmlFor="password">{t('login.lbl_password')}</label>
                <input
                    id="password"
                    type="password"
                    autoComplete="off"
                    placeholder={t('login.ph_password')}
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                />
            </div>

            <label>
                <input
                    type="checkbox"
                    name="rememberMe"
                    checked={rememberMe}
                    onChange={(e) => setRememberMe(e.target.checked)}
                />
                {t('login.remember_me')}
            </label>

            <button type="submit">{t('login.submit')}</button>
        </FormWrapper>
    );
}
