import { useConfig } from '@/providers/ConfigProvider';
export function Footer() {
    const { config } = useConfig();
    return (
        <footer
            className="
                fixed bottom-0 left-0 w-full
                backdrop-blur-md
                px-6 py-3
                text-sm text-center
                border-t border-[var(--border)]
                shadow-[0_-2px_8px_rgba(0,0,0,0.05)]
                z-[1000]
                bg-[var(--surface)]
                text-[var(--text-secondary)]
                transition-all
            ">
            <p className="my-1">
                &copy; {new Date().getFullYear()} {config.AppName}. All rights reserved.
            </p>
            <p className="my-1">Built with ❤️ using React.</p>
        </footer>
    );
}
