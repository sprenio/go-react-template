import { CheckCircle, AlertCircle, AlertTriangle, Info } from 'lucide-react';

const iconsMap = {
  success: CheckCircle,
  error: AlertCircle,
  warning: AlertTriangle,
  info: Info,
};
type MessageBoxProps = {
  message: string;
  type: 'info' | 'success' | 'warning' | 'error';
};
export function MessageBox({ message, type = 'info' }: MessageBoxProps) {
  if (!message) return null;
  const Icon = iconsMap[type] || Info;

  return (
    <div
      className={
        `msg-type type-${type} m-[var(--space-xs)] flex items-center justify-center rounded-[var(--radius-md)] p-[var(--space-sm)] text-[length:var(--text-md)] font-medium`
      }
    >
      <Icon style={{ marginRight: '0.5rem' }} size={20} />
      {message}
    </div>
  );
}
