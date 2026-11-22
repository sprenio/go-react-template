import { useConfig } from '@/providers/ConfigProvider';
import { useTranslation } from 'react-i18next';

type PageTitleProps = {
  page: string;
  description?: string;
  strKey?: string;
  strParams?: Record<string, string>;
};

export function PageTitle({
  page,
  description,
  strKey = 'page_title',
  strParams = {},
}: PageTitleProps) {
  const { config } = useConfig();
  const { t } = useTranslation();
  const appName = config.AppName;
  let title = t(page + '.' + strKey, strParams);
  if (title > '') {
    title += ' :: ';
  }
  title += appName;
  description = description || appName;
  return (
    <>
      <title>{title}</title>
      <meta name="description" content={description} />
      <meta name="og:description" content={description} />
    </>
  );
}
