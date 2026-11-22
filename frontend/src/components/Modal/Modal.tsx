import type { ReactNode } from 'react';

type ModalProps = {
  open: boolean;
  onClose: () => void;
  children: ReactNode;
  closeOnBackdropClick?: boolean;
};

export function Modal({
  open,
  onClose,
  children,
  closeOnBackdropClick = false,
}: ModalProps) {
  if (!open) return null;

  const handleBackdropClick = () => {
    if (closeOnBackdropClick) {
      onClose();
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      {/* Overlay */}
      <div
        className="absolute inset-0 bg-[var(--overlay)] backdrop-blur-sm transition-opacity"
        onClick={handleBackdropClick}
      />

      {/* Modal */}
      <div className="relative z-10 max-h-[calc(100vh-150px)] w-full max-w-lg overflow-y-auto rounded-[var(--radius-lg)] bg-[var(--surface)] p-[var(--space-lg)] text-[var(--text)] shadow-[var(--shadow-lg)] transition-[var(--transition-normal)]">
        {children}
      </div>
    </div>
  );
}
