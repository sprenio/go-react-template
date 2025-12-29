export const Header = ({title, subtitle}:{title: string, subtitle?: string}) => {
    return (
        <header className={`text-center mb-[var(--space-lg)] fade-in`}>
            <h1 className={`text-[length:var(--text-xl)]
    font-semibold
    text-[var(--primary)]
    my-[var(--space-xs)]
    leading-[1.3]`}>
                {title}
            </h1>
            {subtitle && <p className={`m-0 text-[var(--text-secondary)] text-[length:var(--text-md)] leading-[1.4] font-[var(--font-body)]`}>
                {subtitle}
            </p>}
        </header>
    )
}