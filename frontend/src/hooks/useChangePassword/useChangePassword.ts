import {useMutation} from '@tanstack/react-query';
import {api, apiCodes, ApiError} from '@/api';
import type {ApiSuccessResponse} from '@/api';
import type {ChangePasswordData} from './types'

export function useChangePassword() {
    return useMutation<ApiSuccessResponse, ApiError, ChangePasswordData>({
        mutationFn: async (data: ChangePasswordData) => {
            const response = await api.post<undefined>('/password-change/' + data.hash, data);
            if (response.code !== apiCodes.API_Password_Change_Success) {
                throw new ApiError(response.message || 'Error changing password', response.code);
            }
            return response;
        },
    });
}
