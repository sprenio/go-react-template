import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { api } from '@/api';
import {PageTitle} from '@/layouts/PageTitle';
import { Card } from '@/components/Card';
import { GridSmLgWrapper } from '@/layouts/GridSmLgWrapper';
import type {PingResponse} from './types';

function Dashboard() {
  const [message, setMessage] = useState('...');
  const { t } = useTranslation();

  useEffect(() => {
    api
      .get<PingResponse>('/ping')
      .then((resp) => setMessage(resp.data.message))
      .catch(() => setMessage('Error'));
  }, []);

  return (
    <>
      <PageTitle page={'dashboard'} />
      {/* Dwa mniejsze cardy po lewej */}
      <GridSmLgWrapper>
        <div className="grid-column">
          <Card title={t('dashboard.stats_1_title')}>
            <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</p>
          </Card>
          <Card title={t('dashboard.stats_2_title')}>
            <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</p>
          </Card>
        </div>

        {/* Duży card po prawej */}
        <div className="grid-column">
          <Card title={t('dashboard.panel_title')}>
            <p dangerouslySetInnerHTML={{ __html: t('dashboard.backend_response', { message }) }} />
          </Card>
        </div>
      </GridSmLgWrapper>
    </>
  );
}

export default Dashboard;
