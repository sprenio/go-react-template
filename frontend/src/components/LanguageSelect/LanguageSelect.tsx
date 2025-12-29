import { useState, useEffect } from 'react';
import { useConfig } from '@/providers/ConfigProvider';
import { useLang } from '@/providers/LangProvider';
import type {LangType} from '@/types';

type LanguageSelectProps = {
    value: string
    onChange: (newLang:string)=>void
    flagsOnly?:boolean
}
export default function LanguageSelect({ value = '', onChange, flagsOnly = false }: LanguageSelectProps) {
    const { lang: defaultLang } = useLang();
    value = value || defaultLang;
    const { config } = useConfig();
    const languages = config.Languages || [];
    const current = languages.find((lang:LangType) => lang.code === value) || languages[0];
    const [selected, setSelected] = useState(current);
    const [open, setOpen] = useState(false);

    useEffect(() => {
        if (languages.length > 0) {
            const found = languages.find((l) => l.code === value) || languages[0];
            setSelected(found);
        }
    }, [languages, value]);

    const handleSelect = (lang:LangType) => {
        setSelected(lang);
        setOpen(false);
        if (onChange) {
            onChange(lang.code);
        }
    };
    return (
        <div className={`relative text-[length:var(--text-sm)] ${flagsOnly ? 'w-[68px]' : 'w-[200px]'}`}>
            <button className=" w-full flex items-center justify-between bg-[var(--surface)] text-[var(--text)] border border-[var(--border)] rounded-[var(--radius-md)] px-[var(--space-sm)] py-[var(--space-xs)] cursor-pointer transition-colors duration-200 hover:border-[var(--primary)]"
                    type="button" onClick={() => setOpen(!open)}>
                <span className={`flag flag-${selected?.code}`} />
                {flagsOnly || <span>{selected?.name}</span>}
                <span>▾</span>
            </button>

            {open && (
                <div className=" absolute top-full left-0 w-full mt-1 bg-[var(--surface)] border border-[var(--border)] rounded-[var(--radius-md)] shadow--[var(--shadow-sm)] z-10 overflow-hidden">
                    {languages.map((lang) => (
                        <button key={lang.code} type="button" onClick={() => handleSelect(lang)}
                                className="w-full flex items-center bg-transparent border-0 px-[var(--space-sm)] py-[var(--space-xs)] text-left
   text-[var(--text)] cursor-pointer transition-colors duration-200 hover:bg-[var(--surface-hover)]">
                            <span className={`flag flag-${lang.code}`} title={flagsOnly ? lang.name : ''} />
                            {flagsOnly || lang.name}
                        </button>
                    ))}
                </div>
            )}
        </div>
    );
}
