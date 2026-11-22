import { useConfig } from '@/providers/ConfigProvider';
import { Activity } from 'react';
export function Footer() {
  const { config } = useConfig();
  const commit = import.meta.env.VITE_GIT_COMMIT || '';
  const tag = import.meta.env.VITE_GIT_TAG || '_';
  return (
    <footer className="fixed bottom-0 left-0 z-[1000] w-full border-t border-[var(--border)] bg-[var(--surface)] px-6 py-2.5 text-center text-sm text-[var(--text-secondary)] shadow-[0_-2px_8px_rgba(0,0,0,0.05)] backdrop-blur-md transition-all">
      <div className="relative flex items-center justify-center">
        <p className="my-0">
          &copy; {new Date().getFullYear()} {config.AppName}. All rights reserved.
        </p>
        {/* Version */}
        <div className="absolute right-0 flex items-center gap-1 font-mono text-xs text-[var(--text-secondary)] opacity-50 select-none hover:opacity-100 transition-opacity">
          <span>{tag}</span>
          <Activity mode={commit === '' ? 'hidden' : 'visible'} >
            <span>({commit})</span>
          </Activity>
        </div>
      </div>
    </footer>
  );
}
