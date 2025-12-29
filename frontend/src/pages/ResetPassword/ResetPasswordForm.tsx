import React, {useState} from 'react';
import {useTranslation} from 'react-i18next';
import {useNavigate} from 'react-router-dom';
import {getApiCodeDescription} from '@/api';
import {useMessage} from '@/providers/MessageProvider';
import {FormWrapper} from '@/components/Form';
import {useResetPassword} from '@/hooks/useResetPassword';

type ResetFormProps = {
    onSuccess?: () => void
}

function ResetPasswordForm({onSuccess}: ResetFormProps) {
    const {t} = useTranslation();
    const [email, setEmail] = useState('');
    const navigate = useNavigate();
    const {showMessage} = useMessage();
    const resetMutation = useResetPassword();

    const handleSubmit = (e:React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        if(email.trim() === ''){
            return;
        }
        resetMutation.mutate(
            { email },
            {
                onSuccess: () => {
                    onSuccess?.();
                    showMessage(t('reset.success'), 'success');
                    navigate('/login');
                },
                onError: (error) => {
                    let errorMessage = getApiCodeDescription(error.code);
                    showMessage(errorMessage || t('reset.error_unknown'), 'error');
                }
            }
        )
        setEmail('');
    };

    return (
        <FormWrapper onSubmit={handleSubmit}>
            <div>
                <label htmlFor="resetemail">{t('reset.lbl_email')}</label>
                <input
                    type="email"
                    value={email}
                    autoComplete="off"
                    id="resetemail"
                    onChange={(e) => setEmail(e.target.value)}
                    required
                />
            </div>
            <button type="submit" disabled={resetMutation.isPending}>
                {resetMutation.isPending ? t('reset.pending') : t('reset.submit')}
            </button>
        </FormWrapper>
    );
}

export default ResetPasswordForm;
