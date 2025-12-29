import React, {useState} from 'react';
import {useTranslation} from 'react-i18next';
import {useNavigate, useParams} from 'react-router-dom';
import {getApiCodeDescription} from '@/api';
import {toast} from 'react-hot-toast';
import {useMessage} from '@/providers/MessageProvider';
import {FormWrapper} from '@/components/Form';
import {useChangePassword} from '@/hooks/useChangePassword';

type ChangePasswordProps = {
    onSuccess?: () => void
}

function ChangePasswordForm({onSuccess}: ChangePasswordProps) {
    const {hash} = useParams();
    const {t} = useTranslation();
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const changePasswordMutation = useChangePassword();
    const navigate = useNavigate();
    const {showMessage} = useMessage();

    const PASSWORD_REGEX = /^(?=.*[A-Za-z])(?=.*\d)(?=.*[^A-Za-z\d]).{6,}$/;

    const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        if (password !== confirmPassword) {
            showMessage(t('register.error_match'), 'error');
            return;
        }
        if (!PASSWORD_REGEX.test(password)) {
            showMessage(t('register.error_weak'), 'error');
            return;
        }

        changePasswordMutation.mutate(
            {password, newPassword: confirmPassword, hash},
            {
                onSuccess: () => {
                    setPassword('');
                    setConfirmPassword('');
                    onSuccess?.();
                    showMessage(t('password_change.success'), 'success');
                    navigate('/login');
                },
                onError: (error) => {
                    let errorMessage = getApiCodeDescription(error.code);
                    toast.error(errorMessage);
                },
            }
        );
    };

    return (
        <FormWrapper onSubmit={handleSubmit}>
            <div>
                <label htmlFor="password">{t('password_change.lbl_password')}</label>
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
                <label htmlFor="confirm">{t('password_change.lbl_confirm')}</label>
                <input
                    type="password"
                    id="confirm"
                    autoComplete="off"
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    required
                />
            </div>
            <button type="submit">{t('password_change.submit')}</button>
        </FormWrapper>
    );
}

export default ChangePasswordForm;
