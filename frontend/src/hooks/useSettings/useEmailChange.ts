import {useMutation} from '@tanstack/react-query';
import {api, apiCodes, ApiError} from '@/api';
import type {ApiSuccessResponse} from '@/api';
import type {ChangeEmailData} from './types';

export function useEmailChange() {
    return useMutation<ApiSuccessResponse, ApiError, ChangeEmailData>({
        mutationFn: async (data: ChangeEmailData) => {
            const response = await api.post<undefined>('/email_change', data);
            if (response.code !== apiCodes.API_Email_Change_Success) {
                throw new ApiError(response.message || 'Error saving user email', response.code);
            }
            return response;
        },
    });
}
