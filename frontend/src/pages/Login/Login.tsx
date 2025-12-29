import LoginForm from './LoginForm';
import {PageTitle} from '@/layouts/PageTitle';
import { useTranslation } from 'react-i18next';
import { Card } from '@/components/Card';
import { FormFooter, FormHeader } from '@/components/Form';
import { useConfig } from '@/providers/ConfigProvider';
import { Link } from 'react-router-dom';

export function Login() {
    const { t } = useTranslation();
    const { config, loading } = useConfig();
    return (
        <>
            <PageTitle page={'login'} />
            <Card width={360}>
                <FormHeader title={t('login.header')} />
                <LoginForm />
                {!loading && (
                    <FormFooter>
                        {config?.Features.register && (
                            <p>
                                {t('login.no_account')} <Link to="/register">{t('login.register')}</Link>
                            </p>
                        )}
                        {config?.Features.reset_password && (
                            <p>
                                <Link to="/reset-password">{t('login.forgot_password')}</Link>
                            </p>
                        )}
                    </FormFooter>
                )}
            </Card>
        </>
    );
}
