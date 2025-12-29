import { useTranslation } from 'react-i18next';
import {PageTitle} from '@/layouts/PageTitle';
import { Link } from 'react-router-dom';
import { Card } from '@/components/Card';
import { FormFooter, FormHeader } from '@/components/Form';

export function NotFound() {
  const { t } = useTranslation();

  return (
    <>
      <PageTitle page={'not-found'} />
      <Card>
        <FormHeader title={"404"} subtitle={t('not-found.message')} />
        <FormFooter>
          <Link to="/">{t('global.back_to_main_page')}</Link>
        </FormFooter>
      </Card>
    </>
  );
}
