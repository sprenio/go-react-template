import { AlertCircle } from 'lucide-react';

type FormFieldErrorMessageProps = {
    message?: string
};
export function FormFieldErrorMessage({ message }: FormFieldErrorMessageProps) {
    if (!message) return null;

    return (
        <div className='msg-type type-error m-[var(--space-xs)] flex items-center justify-center rounded-[var(--radius-md)] p-[var(--space-sm)] text-[length:var(--text-md)] font-medium'>
            <AlertCircle style={{ marginRight: '0.5rem' }} size={10} />
            {message}
        </div>
    );
}
