import type {SettingsType} from './settings';

export type UserType = {
    name: string
    email: string
    settings?: SettingsType
}