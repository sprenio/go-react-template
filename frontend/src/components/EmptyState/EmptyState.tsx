import type { EmptyStateProps, EmptyStateType } from './types';
import { useTranslation } from 'react-i18next';

const ICONS: Record<EmptyStateType, string> = {
  info: 'ℹ️',
  warning: '⚠️',
  error: '❌',
};

export function EmptyState({ message = '', type = 'info', action }: EmptyStateProps) {
  const { t } = useTranslation();
  message = message || t('global.empty_state_message');
  const icon = ICONS[type];

  return (
    <div
      className={`msg-type flex flex-col items-center justify-center gap-4 rounded-lg border p-6 text-center type-${type}`}
      role="status"
    >
      <div className={`msg-icon flex h-12 w-12 items-center justify-center rounded-full text-xl`}>
        {icon}
      </div>

      <p className={`text-sm`}>{message}</p>

      {action && <div>{action}</div>}
    </div>
  );
}
