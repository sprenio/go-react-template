import { Pencil, Plus, Save, Trash2 } from 'lucide-react';
import { twMerge } from 'tailwind-merge';
import type { ActionButtonProps, ActionButtonsProps, ButtonSize } from './types';
import i18n from '@/i18n';

const baseClasses =
  'p-[var(--space-xs)] rounded-[var(--radius-md)] text-[var(--text-secondary)] hover:text-[var(--primary)] hover:bg-[var(--primary-soft)] transition-colors duration-[var(--transition-fast)]';
const defaultSize = 'md';

export function ActiionButton({ onClick, title, icon, ...props }: ActionButtonProps) {
  const mergedClasses = twMerge(baseClasses, props.className);
  return (
    <button type="button" title={title} onClick={onClick} className={mergedClasses}>
      {icon}
    </button>
  );
}

export function EditButton({ onClick, title, size, ...props }: ActionButtonProps) {
  const buttonSize = getActionButtonSize(size);
  title = title || i18n.t('global.btn_edit');
  return (
    <ActiionButton onClick={onClick} title={title} icon={<Pencil size={buttonSize} />} {...props} />
  );
}

export function DeleteButton({ onClick, title, size, ...props }: ActionButtonProps) {
  const buttonSize = getActionButtonSize(size);
  title = title || i18n.t('global.btn_delete');
  return (
    <ActiionButton
      onClick={onClick}
      title={title}
      icon={
        <Trash2
          size={buttonSize}
          className="hover:bg-[var(--error-bg)] hover:text-[var(--error)]"
        />
      }
      {...props}
    />
  );
}

export function CreateButton({ onClick, title, size, ...props }: ActionButtonProps) {
  const buttonSize = getActionButtonSize(size);
  title = title || i18n.t('global.btn_create');
  return (
    <ActiionButton onClick={onClick} title={title} icon={<Plus size={buttonSize} />} {...props} />
  );
}

export function SaveButton({ onClick, title, size, ...props }: ActionButtonProps) {
  const buttonSize = getActionButtonSize(size);
  title = title || i18n.t('global.btn_save');
  return (
    <ActiionButton onClick={onClick} title={title} icon={<Save size={buttonSize} />} {...props} />
  );
}

export function ActionButtons({
  editButton,
  deleteButton,
  createButton,
  saveButton,
}: ActionButtonsProps) {
  return (
    <div className="flex items-center gap-[var(--space-xs)]">
      {createButton?.show && <CreateButton {...createButton} />}
      {saveButton?.show && <SaveButton {...saveButton} />}
      {editButton?.show && <EditButton {...editButton} />}
      {deleteButton?.show && <DeleteButton {...deleteButton} />}
    </div>
  );
}

function getActionButtonSize(size: ButtonSize | undefined): number {
  size = size || defaultSize;
  switch (size) {
    case 'sm':
      return 12;
    case 'lg':
      return 20;
    case 'md':
    default:
      return 16;
  }
}
