import { CheckCircle, AlertCircle, AlertTriangle, Info } from 'lucide-react';

const iconsMap = {
  success: CheckCircle,
  error: AlertCircle,
  warning: AlertTriangle,
  info: Info,
};
type MessageBoxProps = {
    message: string,
    type: 'info' | 'success' | 'warning' | 'error'
}
export function MessageBox({ message, type = 'info' }:MessageBoxProps) {
  if (!message) return null;

  const Icon = iconsMap[type] || Info;

  return (
    <div className={`flex items-center justify-center p-[var(--space-sm)] m-[var(--space-xs)] rounded-[var(--radius-md)] text-[length:var(--text-md)] font-medium
        msg-type type-` + type}>
      <Icon style={{ marginRight: '0.5rem' }} size={20} />
      {message}
    </div>
  );
}
