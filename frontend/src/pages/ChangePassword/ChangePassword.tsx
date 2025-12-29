import {PageTitle} from '@/layouts/PageTitle';
import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';
import { useConfig } from '@/providers/ConfigProvider';
import {NotFound} from '../NotFound';
import { FormFooter, FormHeader } from '@/components/Form';
import { Card } from '@/components/Card';
import ChangePasswordForm from './ChangePasswordForm';

function ChangePassword() {
  const { config, loading } = useConfig();
  const { t } = useTranslation();

  if (!config.Features?.reset_password) {
    return loading ? '' : <NotFound />;
  }
  return (
    <>
      <PageTitle page="password_change" />
      <Card>
        <FormHeader title={t('password_change.title')} subtitle={t('password_change.subtitle')} />
        <ChangePasswordForm />
        <FormFooter>
          <p>
            <Link to="/login">{t('reset.back_to_login')}</Link>
          </p>
        </FormFooter>
      </Card>
    </>
  );
}

export default ChangePassword;
