import { useTranslation } from 'react-i18next'

type LoaderProps = {
    message?: string
}

export default function Loader({ message }: LoaderProps) {
    const { t } = useTranslation()
    message = message || t('loader.please_wait')
    return (
        <div className="fixed inset-0 z-[9999] flex items-center justify-center">
            {/* overlay */}
            <div className="absolute inset-0 bg-black/40 backdrop-blur-sm" />

            {/* loader box */}
            <div className="relative z-10 flex flex-col items-center gap-4 rounded-xl bg-[var(--surface)] px-8 py-6 shadow-xl">
                {/* spinner */}
                <div className="h-10 w-10 animate-spin rounded-full border-4 border-[var(--border)] border-t-[var(--primary)]" />

                {/* message */}
                <p className="text-sm text-[var(--text-secondary)] text-center">
                    {message}
                </p>
            </div>
        </div>
    )
}
