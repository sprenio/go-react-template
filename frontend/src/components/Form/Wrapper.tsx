import React from 'react';
import { useTranslation } from 'react-i18next';

type Props = {
  children: React.ReactNode;
  onSubmit: (e: React.FormEvent<HTMLFormElement>) => void;
  showRequiredFieldsInfo?: boolean;
};

export const Wrapper = ({ children, onSubmit, showRequiredFieldsInfo = true }: Props) => {
  const { t } = useTranslation();
  return (
    <form className={'form-wrapper'} onSubmit={onSubmit}>
      {showRequiredFieldsInfo && (
        <p
          className={'required-fields-info'}
          dangerouslySetInnerHTML={{ __html: t('global.required_fields_info') }}
        />
      )}
      {children}
    </form>
  );
};
