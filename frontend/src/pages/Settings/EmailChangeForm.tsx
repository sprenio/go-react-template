import { FormWrapper } from '@/components/Form';
import { useTranslation } from 'react-i18next';
import React, { useState } from 'react';


type EmailChangeFormProps = {
    email: string
    onSubmit: (email: string) => void
}
export default function EmailChangeForm({ email, onSubmit }:EmailChangeFormProps) {
  const { t } = useTranslation();
  const [newEmail, setNewEmail] = useState('');

  const handleEmailChange = (e:React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (onSubmit) {
      onSubmit(newEmail);
    }
  };

  return (
    <FormWrapper onSubmit={handleEmailChange}>
      <p>{t('settings.email_change_description')}</p>

      {/* Email */}
      <div>
        <label htmlFor="newEmail">{t('settings.lbl_new_email')}</label>
        <input
          id="newEmail"
          type="email"
          placeholder={email}
          value={newEmail}
          onChange={(e) => setNewEmail(e.target.value)}
        />
      </div>
      <button type="submit">{t('settings.btn_save')}</button>
    </FormWrapper>
  );
}
