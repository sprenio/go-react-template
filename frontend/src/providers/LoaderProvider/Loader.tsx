import { useTranslation } from 'react-i18next';

type LoaderProps = {
  message?: string;
};

export default function Loader({ message }: LoaderProps) {
  const { t } = useTranslation();
  const displayMessage = message || t('loader.please_wait');

  return (
    <div className="fixed inset-0 z-[9999] flex items-center justify-center">
      {/* Tło / overlay */}
      <div className="absolute inset-0 bg-black/50 backdrop-blur-sm" />

      {/* Karta loadera */}
      <div className="relative z-10 flex min-w-[220px] max-w-sm flex-col items-center gap-4 rounded-2xl border border-[var(--border)] bg-[var(--surface)] px-8 py-7 shadow-2xl">
        {/* Kontener animacji spinnera */}
        <div className="relative flex h-14 w-14 items-center justify-center">
          {/* Zewnętrzny pierścień tła */}
          <div className="absolute inset-0 rounded-full border-4 border-[var(--border)] opacity-40" />

          {/* Aktywny obracający się pierścień z akcentem motywu */}
          <div className="absolute inset-0 animate-spin rounded-full border-4 border-transparent border-t-[var(--primary)] border-r-[var(--primary-hover)]" />

          {/* Wewnętrzny pulsujący element akcentujący */}
          <div className="h-3 w-3 animate-ping rounded-full bg-[var(--primary)] opacity-75" />
          <div className="absolute h-2 w-2 rounded-full bg-[var(--primary)]" />
        </div>

        {/* Komunikat ze stylowym przejściem i zmiennymi motywu */}
        <p className="animate-pulse text-center font-[var(--font-body)] text-[length:var(--text-md)] font-[var(--font-medium)] text-[var(--text-secondary)]">
          {displayMessage}
        </p>
      </div>
    </div>
  );
}

