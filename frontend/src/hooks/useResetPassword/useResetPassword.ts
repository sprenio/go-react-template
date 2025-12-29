import { useMutation } from '@tanstack/react-query';
import {api, apiCodes, ApiError} from '@/api';
import type {ApiSuccessResponse} from '@/api';
import type {ResetPasswordData} from './types';

export function useResetPassword() {
    return useMutation<ApiSuccessResponse, ApiError, ResetPasswordData>({
        mutationFn: async (data: ResetPasswordData) => {
            const response = await api.post<undefined>('/reset-password', data);
            if (response.code !== apiCodes.API_Reset_Password_Success) {
                throw new ApiError(response.message || 'Reset password failed', response.code);
            }
            return response;
        },
    });
}
