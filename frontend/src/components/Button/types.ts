import type { ReactNode, ButtonHTMLAttributes } from 'react';

export type ButtonType =
  | 'primary'
  | 'secondary'
  | 'success'
  | 'warning'
  | 'error'
  | 'info'
  | 'link'
  | 'soft';

export type ButtonSize = 'sm' | 'md' | 'lg';
export type IconPosition = 'left' | 'right' | 'none';

export type ButtonProps = Omit<ButtonHTMLAttributes<HTMLButtonElement>, 'type'> & {
  children: ReactNode;

  type?: ButtonType;
  size?: ButtonSize;

  htmlType?: 'button' | 'submit' | 'reset';

  icon?: ReactNode;
  iconPosition?: IconPosition;

  isLoading?: boolean;
};

export type ActionButtonProps = {
  onClick?: () => void;
  title?: string;
  icon?: ReactNode;
  size?: ButtonSize;
} & Omit<ButtonHTMLAttributes<HTMLButtonElement>, 'type'>;

export type ActionButtonsProps = {
  editButton?: ActionButtonProps & { show?: boolean };
  deleteButton?: ActionButtonProps & { show?: boolean };
  createButton?: ActionButtonProps & { show?: boolean };
  saveButton?: ActionButtonProps & { show?: boolean };
};
