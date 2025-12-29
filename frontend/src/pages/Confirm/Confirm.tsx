import {useParams} from 'react-router-dom';
import {useEffect, useRef, useState} from 'react';
import {useTranslation} from 'react-i18next';
import {useNavigate} from 'react-router-dom';
import {NotFound} from '../NotFound';
import {useConfirm} from '@/hooks/useConfirm';
import {apiCodes, ApiError, apiErrorCodes, getApiCodeDescription} from '@/api';
import {useMessage} from '@/providers/MessageProvider';

function Confirm() {
    const {hash} = useParams();
    const {t} = useTranslation();
    const navigate = useNavigate();
    const {showMessage} = useMessage();
    const didFetch = useRef(false);
    const confirmMutation = useConfirm();
    const [isPending, setIsPending] = useState<boolean>(true);

    useEffect(() => {
        if (didFetch.current) return;
        didFetch.current = true;
        if (!hash) {
            return;
        }
        const runConfirm = async() => {
            const messageMap: Record<string, string> = {
                email_change: t('confirm.success_email_change'),
                register: t('confirm.success_register'),
            };
            try {
                const data = await confirmMutation.mutateAsync({hash});
                const message = messageMap[data.token_type] ?? t('confirm.success_message');
                showMessage(message, 'success');
                navigate('/');
            } catch (error:unknown) {
                let errorMessage = getApiCodeDescription(apiCodes.API_General_Unknown_Error);
                if (error instanceof ApiError) {
                    if (error.code == apiErrorCodes.NOT_FOUND) {
                        return;
                    }
                    console.error('Error confirming token:', error.code, error.message, error);
                    errorMessage = getApiCodeDescription(error.code);
                } else {
                    console.error('Unknown error', error);
                }
                showMessage(errorMessage, 'error');
                navigate('/login');
            }
        };
        runConfirm().finally(() => {setIsPending(false)});

    }, [hash]);

    if (isPending) {
        return <div></div>;
    }
    return <NotFound/>;
}

export default Confirm;
