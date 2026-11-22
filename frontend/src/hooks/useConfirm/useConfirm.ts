import { useMutation } from '@tanstack/react-query';
import { api, apiCodes, ApiError } from '@/api';
import type { ConfirmData, ConfirmResponse } from './types';

export function useConfirm() {
  return useMutation<ConfirmResponse, ApiError, ConfirmData>({
    mutationFn: async (data: ConfirmData) => {
      const response = await api.get<ConfirmResponse>('/confirm/' + data.hash);
      if (response.code !== apiCodes.API_Confirm_Success) {
        throw new ApiError(response.message || 'Confirmation failed', response.code);
      }
      return response.data;
    },
  });
}
