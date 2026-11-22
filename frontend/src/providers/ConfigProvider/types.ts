import type { ConfigType } from '@/domain';

export type ConfigContextType = {
  config: ConfigType;
  showLanguages: boolean;
  loading: boolean;
};
export type ApiConfigResponse = {
  data: ConfigType;
};
