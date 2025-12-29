type LogoProps = {
    theme: 'light' | 'dark'
}

export function Logo({ theme }: LogoProps) {
    const isDark = theme === 'dark'
    return (
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 110 60">
            <circle cx="30" cy="30" r="18" fill={isDark ? '#61dafb' : '#0ea5e9'}/>
            <circle cx="80" cy="30" r="18" fill={isDark ? '#22c55e' : '#16a34a'}/>
            <rect x="48" y="27" width="14" height="6" rx="3" fill={isDark ? '#94a3b8' : '#64748b'}/>
            {/*
            <text
                x="110"
                y="38"
                font-family="Inter, system-ui, -apple-system, sans-serif"
                font-size="24"
                font-weight="600"
                fill={isDark ? '#e5e7eb' : '#0f172a'}
                letter-spacing="0.5"
            >
                {appName}
            </text>
        */}
        </svg>
    );
}
