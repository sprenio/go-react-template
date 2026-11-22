import { useMutation } from '@tanstack/react-query';
import { api, apiCodes, ApiError } from '@/api';
import type { ApiSuccessResponse } from '@/api';
import type { SettingsType } from '@/domain';

export function useSettings() {
  return useMutation<ApiSuccessResponse, ApiError, SettingsType>({
    mutationFn: async (data: SettingsType) => {
      const response = await api.post<undefined>('/settings', data);
      if (response.code !== apiCodes.API_General_Success) {
        throw new ApiError(response.message || 'Error saving user settings', response.code);
      }
      return response;
    },
  });
}
