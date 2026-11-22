import { useState, useEffect, useMemo } from 'react';
import { useConfig } from '@/providers/ConfigProvider';
import { useLang } from '@/providers/LangProvider';
import type { LangType } from '@/domain';

type LanguageSelectProps = {
  value: string;
  onChange: (newLang: string) => void;
  flagsOnly?: boolean;
};
export default function LanguageSelect({
  value = '',
  onChange,
  flagsOnly = false,
}: LanguageSelectProps) {
  const { config, showLanguages } = useConfig();

  const { lang: defaultLang } = useLang();
  value = value || defaultLang;
  const languages = useMemo(() => config.Languages ?? [], [config.Languages]);

  const current = languages.find((lang: LangType) => lang.code === value) || languages[0];
  const [selected, setSelected] = useState(current);
  const [open, setOpen] = useState(false);

  useEffect(() => {
    if (languages.length > 0) {
      const found = languages.find((l) => l.code === value) || languages[0];
      setSelected(found);
    }
  }, [languages, value]);

  if (!showLanguages) {
    return null;
  }
  const handleSelect = (lang: LangType) => {
    setSelected(lang);
    setOpen(false);
    if (onChange) {
      onChange(lang.code);
    }
  };
  return (
    <div
      className={`relative text-[length:var(--text-sm)] ${flagsOnly ? 'w-[68px]' : 'w-[200px]'}`}
    >
      <button
        className="flex w-full cursor-pointer items-center justify-between rounded-[var(--radius-md)] border border-[var(--border)] bg-[var(--surface)] px-[var(--space-sm)] py-[var(--space-xs)] text-[var(--text)] transition-colors duration-200 hover:border-[var(--primary)]"
        type="button"
        onClick={() => setOpen(!open)}
      >
        <span className={`flag flag-${selected?.code}`} />
        {flagsOnly || <span>{selected?.name}</span>}
        <span>▾</span>
      </button>

      {open && (
        <div className="shadow--[var(--shadow-sm)] absolute top-full left-0 z-10 mt-1 w-full overflow-hidden rounded-[var(--radius-md)] border border-[var(--border)] bg-[var(--surface)]">
          {languages.map((lang) => (
            <button
              key={lang.code}
              type="button"
              onClick={() => handleSelect(lang)}
              className="flex w-full cursor-pointer items-center border-0 bg-transparent px-[var(--space-sm)] py-[var(--space-xs)] text-left text-[var(--text)] transition-colors duration-200 hover:bg-[var(--surface-hover)]"
            >
              <span className={`flag flag-${lang.code}`} title={flagsOnly ? lang.name : ''} />
              {flagsOnly || lang.name}
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
