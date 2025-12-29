import {useMutation} from '@tanstack/react-query';
import {api, apiCodes, ApiError} from '@/api';
import type {ApiSuccessResponse} from '@/api';
import type {RegisterData} from './types';

export function useRegister() {
    return useMutation<ApiSuccessResponse, ApiError, RegisterData>({
        mutationFn: async (data: RegisterData) => {
            const response = await api.post<undefined>('/register', data);
            if (response.code !== apiCodes.API_Register_Success) {
                throw new ApiError(response.message || 'Register failed', response.code);
            }
            return response;
        },
    });
}
