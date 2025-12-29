import { useTranslation } from 'react-i18next'

export default function Loader({ message = "" }) {
    const { t } = useTranslation()
    message = message || t('loader.please_wait')
    return (
        <div>
            <div></div>
            <p>{message}</p>
        </div>
    )
}
