import {PageTitle} from '@/layouts/PageTitle';
import { useTranslation } from 'react-i18next';
import ResetPasswordForm from './ResetPasswordForm';
import { Link } from 'react-router-dom';
import { useConfig } from '@/providers/ConfigProvider';
import {NotFound} from '../NotFound';
import { FormFooter, FormHeader } from '@/components/Form';
import { Card } from '@/components/Card';

function ResetPassword() {
  const { config, loading } = useConfig();
  const { t } = useTranslation();

  if (!config.Features?.reset_password) {
    return loading ? '' : <NotFound />;
  }
  return (
    <>
      <PageTitle page={'reset'} />
      <Card>
        <FormHeader title={t('reset.header')} subtitle={t('reset.subtitle')} />
        <ResetPasswordForm />
        <FormFooter>
          <p>
            <Link to="/login">{t('reset.back_to_login')}</Link>
          </p>
        </FormFooter>
      </Card>
    </>
  );
}

export default ResetPassword;
