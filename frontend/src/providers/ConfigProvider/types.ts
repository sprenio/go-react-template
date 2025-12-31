import type {ConfigType} from '@/types';

export type ConfigContextType = {
    config: ConfigType,
    showLanguages:boolean,
    loading: boolean
}
export type ApiConfigResponse = {
    data: ConfigType
}