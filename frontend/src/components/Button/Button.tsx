import type { ReactNode } from 'react';
import type { ButtonProps, ButtonType, ButtonSize } from './types';

const baseClasses =
  'inline-flex items-center gap-2 rounded-md font-medium transition focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:opacity-50 disabled:pointer-events-none';

const sizeClasses: Record<ButtonSize, string> = {
  sm: 'px-[var(--space-sm)] py-[var(--space-xs)] text-[length:var(--text-sm)]',
  md: 'px-[var(--space-md)] py-[var(--space-sm)] text-[length:var(--text-md)]',
  lg: 'px-[var(--space-lg)] py-[var(--space-md)] text-[length:var(--text-lg)]',
};

const typeClasses: Record<ButtonType, string> = {
  primary:
    'bg-[var(--primary)] text-[var(--text-on-primary)] hover:bg-[var(--primary-hover)] focus:ring-[var(--primary)]',

  secondary:
    'bg-[var(--secondary)] text-[var(--text-on-secondary)] hover:bg-[var(--secondary-hover)] focus:ring-[var(--secondary)]',

  success: 'bg-[var(--success)] text-white hover:opacity-90 focus:ring-[var(--success)]',

  warning: 'bg-[var(--warning)] text-black hover:opacity-90 focus:ring-[var(--warning)]',

  error: 'bg-[var(--error)] text-white hover:opacity-90 focus:ring-[var(--error)]',

  info: 'bg-[var(--info)] text-white hover:opacity-90 focus:ring-[var(--info)]',

  soft: 'bg-[var(--primary-soft)] text-[var(--text-on-primary-soft)] hover:opacity-90 focus:ring-[var(--primary)]',

  link: 'bg-transparent text-[var(--primary)] px-0 py-0 hover:underline focus:ring-0',
};

const defaultIcons: Partial<Record<ButtonType, ReactNode>> = {
  primary: '⭐',
  secondary: '➡️',
  success: '✅',
  warning: '⚠️',
  error: '❌',
  info: 'ℹ️',
  link: '🔗',
  soft: '✨',
};

const Spinner = () => (
  <span
    className="inline-block h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent"
    aria-hidden
  />
);

const cx = (...classes: (string | undefined | false)[]) => classes.filter(Boolean).join(' ');

export function Button({
  children,
  type = 'secondary',
  size = 'md',

  htmlType = 'button',

  icon,
  iconPosition = 'none',

  isLoading = false,
  disabled = false,

  className,
  onClick,

  ...rest
}: ButtonProps) {
  const resolvedIcon = icon ?? defaultIcons[type];
  const showIcon = iconPosition !== 'none' && (resolvedIcon || isLoading);

  return (
    <button
      type={htmlType}
      onClick={onClick}
      disabled={disabled || isLoading}
      className={cx(baseClasses, sizeClasses[size], typeClasses[type], className)}
      {...rest}
    >
      {showIcon && iconPosition === 'left' && (
        <span className="inline-flex">{isLoading ? <Spinner /> : resolvedIcon}</span>
      )}

      <span>{children}</span>

      {showIcon && iconPosition === 'right' && (
        <span className="inline-flex">{isLoading ? <Spinner /> : resolvedIcon}</span>
      )}
    </button>
  );
}
