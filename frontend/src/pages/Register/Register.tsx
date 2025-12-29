import {PageTitle} from '@/layouts/PageTitle';
import RegisterForm from './RegisterForm';
import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';
import { useConfig } from '@/providers/ConfigProvider';
import {NotFound} from '../NotFound';
import { Card } from '@/components/Card';
import { FormFooter, FormHeader } from '@/components/Form';

function Register() {
  const { config, loading } = useConfig();
  const { t } = useTranslation();

  if (!config?.Features.register) {
    return loading ? '' : <NotFound />;
  }

  return (
    <>
      <PageTitle page={'register'} />
      <Card>
        <FormHeader title={t('register.title')} subtitle={t('register.subtitle')} />
        <RegisterForm />
        <FormFooter>
          <p>
            {t('register.have_account')} <Link to="/login">{t('register.login_link')}</Link>
          </p>
        </FormFooter>
      </Card>
    </>
  );
}
export default Register;
