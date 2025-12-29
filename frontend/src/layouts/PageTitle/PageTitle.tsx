import { useConfig} from "@/providers/ConfigProvider";
import {useTranslation} from 'react-i18next'

type PageTitleProps =  {
    page: string,
    description?: string,
}

export function PageTitle({ page, description }: PageTitleProps) {
    const {config} = useConfig();
    const {t} = useTranslation()
    const appName = config.AppName;
    let title = t(page + '.page_title')
    if(title > '') {
        title += ' :: ';
    }
    title += appName;
    description = description || appName;
    return (
        <>
            <title>{title}</title>
            <meta name="description" content={description} />
            <meta name="og:description" content={description} />
        </>
    );
}