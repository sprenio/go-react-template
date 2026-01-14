import {useMutation} from '@tanstack/react-query';
import {api, ApiError} from '@/api';
import i18n from '@/i18n';
import type {MeResponse} from '@/providers/AuthProvider';
import type {UserCredentials} from './types';

async function login(credentials: UserCredentials): Promise<MeResponse> {
    return await api.post<MeResponse>('/login', credentials, i18n.t('login.loader')).then(({data}) => {
        return data
    })
}

export function useLogin() {
    return useMutation<MeResponse, ApiError, UserCredentials>({
        mutationFn: login,
    });
}
