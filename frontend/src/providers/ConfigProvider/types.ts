import type {ConfigType} from '@/types';

export type ConfigContextType = {
    config: ConfigType,
    loading: boolean
}
export type ApiConfigResponse = {
    data: ConfigType
}