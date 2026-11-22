export const Header = ({
  title,
  subtitle,
  description,
}: {
  title: string;
  subtitle?: string;
  description?: string;
}) => {
  return (
    <header className={`fade-in mb-[var(--space-lg)]`}>
      <h1
        className={`my-[var(--space-xs)] text-center text-[length:var(--text-xl)] leading-[1.3] font-semibold text-[var(--primary)]`}
      >
        {title}
      </h1>
      {subtitle && (
        <p
          className={`m-0 text-center text-[length:var(--text-md)] leading-[1.4] font-[var(--font-body)] text-[var(--text-secondary)]`}
        >
          {subtitle}
        </p>
      )}
      {description && (
        <p
          className={`m-0 text-[length:var(--text-sm)] leading-[1.4] font-[var(--font-body)] text-[var(--text-secondary)]`}
          dangerouslySetInnerHTML={{ __html: description }}
        />
      )}
    </header>
  );
};
