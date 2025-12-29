import { useMutation } from '@tanstack/react-query';
import {api, ApiError} from '@/api';
import { toast } from 'react-hot-toast';
import i18n from '@/i18n';
import type {MeResponse} from '@/providers/AuthProvider';
import type {UserCredentials} from './types';

async function login(credentials: UserCredentials): Promise<MeResponse> {
    const resp = await api.post<MeResponse>('/login', credentials, i18n.t('login.loader'));
    if (resp?.token) {
        toast.success(i18n.t('login.success'));
    } else {
        toast.error(i18n.t('global.unknown_error'));
    }
    return resp.data;
}

export function useLogin() {
    return useMutation<MeResponse, ApiError, UserCredentials>({
        mutationFn: login,
    });
}
